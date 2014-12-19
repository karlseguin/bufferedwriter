package bufferedwriter

import (
	. "github.com/karlseguin/expect"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type WorkerTests struct{}

func (_ WorkerTests) Worker(t *testing.T) {
	Expectify(new(WorkerTests), t)
}

func (_ WorkerTests) GeneratesTheCorrectPaths() {
	// test both trailing slash and not
	for _, path := range []string{"/home/test", "/home/test/"} {
		w := NewWorker(12, nil, Configure().Path(path))
		Expect(w.fileTemp).To.Equal(os.TempDir() + "12.tmp")
		Expect(w.fileRoot).To.Equal("/home/test/12_")
	}
}

func (_ WorkerTests) GeneratesTheCorrectPathsWithPrefix() {
	// test both trailing slash and not
	for _, path := range []string{"/home/test", "/home/test/"} {
		w := NewWorker(12, nil, Configure().Path(path).Prefix("bw_"))
		Expect(w.fileTemp).To.Equal(os.TempDir() + "bw_12.tmp")
		Expect(w.fileRoot).To.Equal("/home/test/bw_12_")
	}
}

func (_ WorkerTests) GeneratesTheCorrectPathWhenTrailingSlashIsMissing() {
	w := NewWorker(12, nil, Configure().Path("/home/test"))
	Expect(w.fileTemp).To.Equal(os.TempDir() + "12.tmp")
	Expect(w.fileRoot).To.Equal("/home/test/bw_12_")
}

func (_ WorkerTests) BuffersWritesInMemory() {
	expected := []byte("There ain't no such thing as a free lunch")
	w := NewWorker(1, nil, testConfig(100))
	w.process(expected)
	Expect(w.data).To.Equal(expected)
	assertNoIO()
}

func (_ WorkerTests) WriteExactSize() {
	defer cleanup()
	expected := "There ain't no such thing as a free lunch"
	w := NewWorker(1, nil, testConfig(len([]byte(expected))+32))
	w.process([]byte(expected))
	w.process([]byte("next"))
	assertFile(12, 22, expected, ".log")
	Expect(w.length).To.Equal(6)
}

func (_ WorkerTests) HandleMultipleFlushes() {
	defer cleanup()
	w := NewWorker(1, nil, testConfig(5+32))
	w.process([]byte("aaaa"))
	w.process([]byte("bbbbb"))
	w.process([]byte("cccc"))
	files := testFiles(".log")
	assertContent(files[0], 12, 22, "aaaa")
	assertContent(files[1], 12, 22, "bbbbb")
	Expect(w.length).To.Equal(6)
}

func testConfig(size int) *Configuration {
	return Configure().Size(size).Path("/tmp").Prefix("recorder_")
}

func assertNoIO() {
	files := testFiles("*")
	Expect(len(files)).To.Equal(0)
}

func assertFile(time byte, userId byte, expected string, extension string) {
	tmp := testFiles(extension)
	if len(tmp) != 1 {
		Fail("Expecting 1 %v file, got %d", extension, len(tmp))
	} else {
		assertContent(tmp[0], time, userId, expected)
	}
}

func assertContent(file string, time byte, userId byte, expected string) {
	data, _ := ioutil.ReadFile("/tmp/" + file)
	Expect(string(data)).To.Equal(expected)
}

func cleanup() {
	for _, name := range testFiles("*") {
		os.Remove("/tmp/" + name)
	}
}

func testFiles(extension string) []string {
	var matches []string
	files, _ := ioutil.ReadDir("/tmp")
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "recorder_") && (extension == "*" || strings.HasSuffix(name, extension)) {
			matches = append(matches, name)
		}
	}
	return matches
}
