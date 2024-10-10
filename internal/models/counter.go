package models

import "sync"

type Counters struct {
	mx sync.RWMutex
	m  map[string]int
}

func NewCounters() *Counters {
	return &Counters{
		m: make(map[string]int),
	}
}

func (c *Counters) Load(key string) (int, bool) {
	c.mx.RLock()
	val, ok := c.m[key]
	c.mx.RUnlock()
	return val, ok
}

func (c *Counters) Store(key string, value int) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}

func (c *Counters) Inc(key string) {
	c.mx.Lock()
	c.m[key]++
	c.mx.Unlock()
}
