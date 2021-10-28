package core

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type Scheduler interface {
	// Push push a *Context
	Push(c *Context)
	// Pop get a *Context
	// return nil if empty
	Pop() (c *Context)
	// Close close queue
	Close()
}

type LifoScheduler struct {
	closed *int64
	l      *list.List
	lock   sync.Mutex
}

func NewLifoScheduler() *LifoScheduler {
	return &LifoScheduler{
		closed: new(int64),
		l:      list.New(),
		lock:   sync.Mutex{},
	}
}

func (l *LifoScheduler) Push(c *Context) {
	if atomic.LoadInt64(l.closed) != 0 {
		return
	}
	l.lock.Lock()
	l.l.PushBack(c)
	l.lock.Unlock()
}

func (l *LifoScheduler) Pop() (c *Context) {
	l.lock.Lock()
	defer l.lock.Unlock()
	front := l.l.Front()
	if front != nil {
		c = l.l.Remove(front).(*Context)
	}
	return
}

func (l *LifoScheduler) Close() {
	atomic.StoreInt64(l.closed, 1)
}
