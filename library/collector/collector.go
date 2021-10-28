package collector

import (
	"sync"
	"sync/atomic"
)

// Collector collection count.
type Collector struct {
	count int64
}

// NewCollect create a Collector instance
func NewCollect(i int64) *Collector {
	c := &Collector{count: i}
	return c
}

// Add adds delta to the Collector count.
func (c *Collector) Add(delta int64) {
	atomic.AddInt64(&c.count, delta)
}

// Count return Collector count.
func (c *Collector) Count() int64 {
	return c.count
}

// GroupCollector is a set of Collector
type GroupCollector struct {
	m sync.Map
}

// NewGroupCollect create a GroupCollector instance.
func NewGroupCollect() *GroupCollector {
	return &GroupCollector{
		m: sync.Map{},
	}
}

// Add adds delta to the Collector count.
func (c *GroupCollector) Add(key string, delta int64) {
	v, ok := c.m.LoadOrStore(key, NewCollect(delta))
	if !ok {
		return
	}
	v.(*Collector).Add(delta)
}

// Range wrap sync.Map Range method
func (c *GroupCollector) Range(f func(string, *Collector) bool) {
	c.m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Collector))
	})
}
