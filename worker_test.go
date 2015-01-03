package bufferedwriter

import (
	. "github.com/karlseguin/expect"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type WorkerTests struct{}

func Test_Worker(t *testing.T) {
	Expectify(new(WorkerTests), t)
}

func (_ WorkerTests) GeneratesTheCorrectPaths() {
	pid := strconv.Itoa(os.Getpid())
	for _, path := range []string{"/home/test", "/home/test/"} {
		w := NewWorker(12, nil, Configure().Path(path))
		Expect(w.fileTemp).To.Equal(os.TempDir() + "12_" + pid + ".tmp")
		Expect(w.fileRoot).To.Equal("/home/test/12_" + pid + "_")
	}
}

func (_ WorkerTests) GeneratesTheCorrectPathsWithPrefix() {
	pid := strconv.Itoa(os.Getpid())
	for _, path := range []string{"/home/test", "/home/test/"} {
		w := NewWorker(12, nil, Configure().Path(path).Prefix("bw_"))
		Expect(w.fileTemp).To.Equal(os.TempDir() + "bw_12_" + pid + ".tmp")
		Expect(w.fileRoot).To.Equal("/home/test/bw_12_" + pid + "_")
	}
}

func (_ WorkerTests) BuffersWritesInMemory() {
	expected := "There ain't no such thing as a free lunch"
	w := NewWorker(1, nil, testConfig(100))
	w.process(closer(expected))
	Expect(w.data[:len(expected)]).To.Equal([]byte(expected))
	assertNoIO()
}

func (_ WorkerTests) WriteExactSize() {
	expected := "There ain't no such thing as a free lunch"
	w := NewWorker(1, nil, testConfig(len([]byte(expected))))
	w.process(closer(expected))
	w.process(closer("next"))
	assertFile(expected, ".log")
	Expect(w.length).To.Equal(4)
}

func (_ WorkerTests) HandleMultipleFlushes() {
	w := NewWorker(1, nil, testConfig(5))
	w.process(closer("aaaa"))
	w.process(closer("bbbbb"))
	w.process(closer("cccc"))
	files := testFiles(".log")
	assertContent(files[0], "aaaa")
	assertContent(files[1], "bbbbb")
	Expect(w.length).To.Equal(4)
}

func (_ WorkerTests) ForcesAFlush() {
	c := make(chan BytesCloser)
	w := NewWorker(1, c, testConfig(10).Forced(time.Millisecond*5))
	go w.work()
	c <- closer("aa123")
	time.Sleep(time.Millisecond * 10)
	files := testFiles(".log")
	assertContent(files[0], "aa123")
	Expect(w.length).To.Equal(0)
}

func (_ WorkerTests) Each(test func()) {
	os.Mkdir("/tmp/bw/", 0700)
	test()
	os.RemoveAll("/tmp/bw/")
}

func testConfig(size int) *Configuration {
	return Configure().Size(size).Path("/tmp/bw/").Prefix("recorder_")
}

func assertNoIO() {
	files := testFiles("recorder_*")
	Expect(len(files)).To.Equal(0)
}

func assertFile(expected string, extension string) {
	tmp := testFiles(extension)
	if len(tmp) != 1 {
		Fail("Expecting 1 %v file, got %d", extension, len(tmp))
	} else {
		assertContent(tmp[0], expected)
	}
}

func assertContent(file string, expected string) {
	data, _ := ioutil.ReadFile("/tmp/bw/" + file)
	Expect(string(data)).To.Equal(expected)
}

func testFiles(extension string) []string {
	var matches []string
	files, _ := ioutil.ReadDir("/tmp/bw/")
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "recorder_") && (extension == "*" || strings.HasSuffix(name, extension)) {
			matches = append(matches, name)
		}
	}
	return matches
}

func closer(data string) BytesCloser {
	return &BC{bytes: []byte(data)}
}

type BC struct {
	bytes []byte
}

func (b *BC) Bytes() []byte {
	return b.bytes
}

func (b *BC) Close() error {
	return nil
}
