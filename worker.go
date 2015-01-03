package bufferedwriter

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	sync.Mutex
	length     int
	data       []byte
	capacity   int
	fileRoot   string
	fileTemp   string
	permission os.FileMode
	channel    chan BytesCloser
	flush      int32
	forced     time.Duration
	timer      *time.Timer
}

func NewWorker(id int, channel chan BytesCloser, config *Configuration) *Worker {
	idString := strconv.Itoa(id)
	w := &Worker{
		channel:  channel,
		capacity: config.size,
		data:     make([]byte, config.size),
		fileRoot: config.path,
		fileTemp: config.temp,
		forced:   config.forced,
	}

	if w.forced != 0 {
		w.timer = time.NewTimer(config.forced)
	}

	if w.fileRoot[len(w.fileRoot)-1:] != "/" {
		w.fileRoot += "/"
	}

	if w.fileTemp[len(w.fileTemp)-1:] != "/" {
		w.fileTemp += "/"
	}
	pid := strconv.Itoa(os.Getpid())
	w.fileRoot += config.prefix + idString + "_" + pid + "_"
	w.fileTemp += config.prefix + idString + "_" + pid + ".tmp"
	return w
}

func (w *Worker) work() {
	os.Remove(w.fileTemp)
	if w.timer == nil {
		for {
			w.process(<-w.channel)
		}
	}

	for {
		select {
		case message := <-w.channel:
			w.process(message)
		case <-w.timer.C:
			println("flushing")
			w.Lock()
			w.save()
			w.Unlock()
		}
	}
}

func (w *Worker) process(message BytesCloser) {
	defer message.Close()
	w.Lock()
	defer w.Unlock()

	data := message.Bytes()

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

func (w *Worker) Flush(wg *sync.WaitGroup) {
	w.Lock()
	defer w.Unlock()
	defer wg.Done()
	w.save()
}

func (w *Worker) save() {
	if w.length == 0 {
		w.resetTimer()
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
	w.resetTimer()
}

func (w *Worker) resetTimer() {
	if w.timer == nil {
		return
	}
	w.timer.Reset(w.forced)
}
