package bufferedwriter

import (
	"errors"
	"sync"
)

var (
	TimeoutErr = errors.New("timeout")
)

type BytesCloser interface {
	Bytes() []byte
	Close() error
}

type Master struct {
	*Configuration
	channel chan BytesCloser
	Workers []*Worker
}

func New(config *Configuration) *Master {
	channel := make(chan BytesCloser, config.queueSize)
	master := &Master{
		channel:       channel,
		Configuration: config,
		Workers:       make([]*Worker, config.workers),
	}

	for i := 0; i < config.workers; i++ {
		worker := NewWorker(i, channel, config)
		go worker.work()
		master.Workers[i] = worker
	}

	return master
}

func (m *Master) Write(message BytesCloser) error {
	select {
	case m.channel <- message:
		return nil
	default:
		return TimeoutErr
	}
}

func (m *Master) Flush() {
	wg := new(sync.WaitGroup)
	wg.Add(len(m.Workers))
	for _, w := range m.Workers {
		go w.Flush(wg)
	}
	wg.Wait()
}
