package bufferedwriter

import (
	"errors"
)

var (
	TimeoutErr = errors.New("timeout")
)

type Master struct {
	*Configuration
	channel chan []byte
	Workers []*Worker
}

func New(config *Configuration) *Master {
	channel := make(chan []byte, config.queueSize)
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

func (m *Master) Write(data []byte) (int, error) {
	select {
	case m.channel <- data:
		return len(data), nil
	default:
		return 0, TimeoutErr
	}
}

func (m *Master) Flush() {
	for _, w := range m.Workers {
		w.Save()
	}
}
