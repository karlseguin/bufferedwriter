package bufferedwriter

import (
	"os"
	"time"
)

type Configuration struct {
	size       int
	workers    int
	queueSize  int
	temp       string
	path       string
	prefix     string
	timeout    time.Duration
	permission os.FileMode
}

func Configure() *Configuration {
	return &Configuration{
		workers:    8,
		size:       8192,
		permission: 0400,
		queueSize:  1024,
		path:       os.TempDir(),
		temp:       os.TempDir(),
		timeout:    time.Millisecond * 5,
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

// How long to wait for a worker before dropping a message (5ms)
func (c *Configuration) Timeout(timeout time.Duration) *Configuration {
	c.timeout = timeout
	return c
}

// The permission of the final file (0400)
func (c *Configuration) Permission(permission os.FileMode) *Configuration {
	c.permission = permission
	return c
}
