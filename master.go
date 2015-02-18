package bufferedwriter

import (
	"errors"
	"sync"
)

var (
	TimeoutErr = errors.New("timeout")
)

type Byter interface {
	Bytes() []byte
}

type Master struct {
	*Configuration
	channel chan Byter
	Workers []*Worker
}

func New(config *Configuration) *Master {
	channel := make(chan Byter, config.queueSize)
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

func (m *Master) Write(message Byter) error {
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
