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
		workers:    4,
		size:       8192,
		permission: 0400,
		queueSize:  1024,
		path:       os.TempDir(),
		temp:       os.TempDir(),
		timeout:    time.Millisecond * 5,
	}
}

func (c *Configuration) Size(size int) *Configuration {
	c.size = size
	return c
}

func (c *Configuration) QueueSize(size int) *Configuration {
	c.queueSize = size
	return c
}

func (c *Configuration) Workers(count int) *Configuration {
	c.workers = count
	return c
}

func (c *Configuration) Path(path string) *Configuration {
	c.path = path
	return c
}

func (c *Configuration) Temp(temp string) *Configuration {
	c.temp = temp
	return c
}

func (c *Configuration) Prefix(prefix string) *Configuration {
	c.prefix = prefix
	return c
}

func (c *Configuration) Timeout(timeout time.Duration) *Configuration {
	c.timeout = timeout
	return c
}

func (c *Configuration) Permission(permission os.FileMode) *Configuration {
	c.permission = permission
	return c
}
