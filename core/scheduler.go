package core

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type Scheduler interface {
	// Push push a *Context
	Push(c *Context)
	// BlockPop block pop
	// return nil if closed
	BlockPop() (c *Context)
	// Close close queue
	Close()
}

type LifoScheduler struct {
	c      chan *Context
	closed *int64
	l      *list.List
	lock   sync.Mutex
}

func NewLifoScheduler() *LifoScheduler {
	return &LifoScheduler{
		c:      make(chan *Context, 2<<10),
		closed: new(int64),
		l:      list.New(),
		lock:   sync.Mutex{},
	}
}

func (l *LifoScheduler) Push(c *Context) {
	if atomic.LoadInt64(l.closed) != 0 {
		return
	}
	select {
	case l.c <- c:
	default:
		l.lock.Lock()
		l.l.PushBack(c)
		l.lock.Unlock()
	}
}

func (l *LifoScheduler) BlockPop() (c *Context) {
	select {
	case c = <-l.c:
		return c
	default:
	}
	l.lock.Lock()
	front := l.l.Front()
	if front != nil {
		c = l.l.Remove(front).(*Context)
	}
	l.lock.Unlock()
	if c != nil {
		return
	}
	c = <-l.c
	return c
}

func (l *LifoScheduler) Close() {
	atomic.StoreInt64(l.closed, 1)
	close(l.c)
}
