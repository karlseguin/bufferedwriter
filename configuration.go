package bufferedwriter

import (
	"os"
)

type Configuration struct {
	size      int
	workers   int
	queueSize int
	temp      string
	path      string
	prefix    string
}

func Configure() *Configuration {
	return &Configuration{
		workers:    8,
		size:       8192,
		queueSize:  1024,
		path:       os.TempDir(),
		temp:       os.TempDir(),
	}
}

// The amount of data to hold in memory per worker (8192)
func (c *Configuration) Size(size int) *Configuration {
	c.size = size
	return c
}

// The size of the internal queue should all workers be busy (1024)
func (c *Configuration) QueueSize(size int) *Configuration {
	c.queueSize = size
	return c
}

// The number of workers (8)
func (c *Configuration) Workers(count int) *Configuration {
	c.workers = count
	return c
}

// The directory to store files (os.TempDir())
func (c *Configuration) Path(path string) *Configuration {
	c.path = path
	return c
}

// A directory to store temporary files (os.TempDir())
func (c *Configuration) Temp(temp string) *Configuration {
	c.temp = temp
	return c
}

// A prefix to add to file names ("")
func (c *Configuration) Prefix(prefix string) *Configuration {
	c.prefix = prefix
	return c
}
