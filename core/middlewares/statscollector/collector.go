package statscollector

import (
	"sync"
	"sync/atomic"
)

// Collector count收集器
type Collector struct {
	i int64
}

func NewCollect(i int64) *Collector {
	c := &Collector{i: i}
	return c
}

func (c *Collector) Add(i int64) {
	atomic.AddInt64(&c.i, i)
}

// GroupCollector 组合的状态收集器
type GroupCollector struct {
	sync.Map
}

func NewGroupCollect() *GroupCollector {
	return &GroupCollector{
		Map: sync.Map{},
	}
}

func (c *GroupCollector) Add(key string, i int64) {
	v, ok := c.LoadOrStore(key, NewCollect(i))
	if !ok {
		return
	}
	v.(*Collector).Add(i)
}
