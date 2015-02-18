# Buffered Writer
BufferedWriter is meant for writing large amounts of data to disk with as little performance penalty to the caller. It achieves this by sacrificing message ordering and potentially losing data during abnormal termination.

Given these properties, BufferedWriter is well-suited for log data that is non-critical and can be independently sorted (or where ordering doesn't matter).

Latency is reduced by:

* Having multiple independent workers (each with its own memory and file). One worker flushing to disk doesn't block the other workers
* Writes are buffered (using a buffered channel), so that even if all workers are busy, the caller won't block
* Writes are dropped when no worker is available and the buffer is full. This can be handled by the client

## Usage

```go
// set up a global writer, it's thread safe
writer := bufferedwriter.New(bufferedwriter.Configure())


// write to it
buffer := bytes.NewBufferString("it's over 9000!")
writer.Write(buffer)
```

The `Write` method expects a value that exposes a `Bytes() []byte`, such as `bytes.Buffer`.

Optionally, the supplied buffer can also implement the `io.Closer` interface (`Close() error`). When it does, bufferedwriter will call Close ones the buffer is no longer needed. This is useful when paired with some type of pooled bytes object, such as [BytePool](https://github.com/karlseguin/bytepool).

`Write` returns an error. The only meaningful error is `bufferedwriter.TimeoutErr` which happens when the message could not even be written to memory.

## Configuration
The writer is configured via a chained configuration object. The above example initiated the write with a default configuration. Options, with their defaults, are:

```go
config := bufferedwriter.Configure().
            // bytes to hold in memory, per worker, before flushing to disk
            Size(8192).
            // size of the buffered channel
            QueueSize(1024).
            // number of workers to run
            Workers(8).
            // directory to store files
            Path(os.TempDir()).
            // directory for temporary files
            Temp(os.TempDir()).
            // prefix to append to file names
            Prefix("").
            // ensure a flush to disk at the specified interval
            Forced(time.Duration(0))

Files are originally written to the location specified by `Temp()` then renamed. On Linux systems this means that, if `Temp` and `Path` are on the same filesystem, you get good atomicity guarantees. In other words, once the file appears in `Path`, it'll be a complete file.
