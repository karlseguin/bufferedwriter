package bufferedwriter

import (
	"log"
	"os"
	"strconv"
	"time"
	"sync"
)

type Worker struct {
	sync.Mutex
	length     int
	data       []byte
	capacity   int
	fileRoot   string
	fileTemp   string
	permission os.FileMode
	channel    chan []byte
	flush      int32
}

func NewWorker(id int, channel chan []byte, config *Configuration) *Worker {
	idString := strconv.Itoa(id)
	w := &Worker{
		channel:    channel,
		capacity:   config.size,
		data:       make([]byte, config.size),
		fileRoot:   config.path,
		fileTemp:   config.temp,
	}

	if w.fileRoot[len(w.fileRoot)-1:] != "/" {
		w.fileRoot += "/"
	}

	if w.fileTemp[len(w.fileTemp)-1:] != "/" {
		w.fileTemp += "/"
	}

	w.fileRoot += config.prefix + idString + "_"
	w.fileTemp += config.prefix + idString + ".tmp"
	return w
}

func (w *Worker) work() {
	os.Remove(w.fileTemp)
	for {
		w.process(<-w.channel)
	}
}

func (w *Worker) process(data []byte) {
	w.Lock()
	defer w.Unlock()

	l := len(data)
	if l > w.capacity {
		log.Println("bufferedwriter dropped large message", l)
		return
	}
	if l > w.capacity-w.length {
		w.save()
	}
	n := copy(w.data[w.length:], data)
	if n != l {
		// can this happen?
		log.Println("bufferedwriter faile to copy full message")
		return
	}
	w.length += n
}

func (w *Worker) Flush() {
	w.Lock()
	defer w.Unlock()
	w.save()
}

func (w *Worker) save() {
	if w.length == 0 {
		return
	}
	defer func() { w.length = 0 }()
	f, err := os.OpenFile(w.fileTemp, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Println("bufferedwriter failed to create temp file", err)
		return
	}
	defer f.Close()
	f.Write(w.data[:w.length])

	target := w.fileRoot + strconv.FormatInt(time.Now().UnixNano(), 10) + ".log"
	if err := os.Rename(w.fileTemp, target); err != nil {
		log.Printf("bufferedwriter failed to rename %v to %v. Error: %v\n", w.fileTemp, target, err)
		return
	}
}
