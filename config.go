package logger

import "sync"

type config struct {
	level Level
	lock  sync.RWMutex
}

func newConfig(level Level) *config {
	return &config{level: level}
}

func (c *config) getLevel() Level {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.level
}

func (c *config) setLevel(level Level) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.level = level
}
