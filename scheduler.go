package mellivora

import (
	"sync/atomic"

	"github.com/smallnest/queue"
)

type Scheduler interface {
	// Push push a *Context
	Push([]byte)
	// Pop get a *Context
	// return nil if empty
	Pop() []byte
	// Close close queue
	Close()
}

type LifoScheduler struct {
	closed *int64
	q      queue.Queue
}

func NewLifoScheduler() *LifoScheduler {
	return &LifoScheduler{
		closed: new(int64),
		q:      queue.NewSliceQueue(0),
	}
}

func (l *LifoScheduler) Push(c []byte) {
	if atomic.LoadInt64(l.closed) != 1 {
		l.q.Enqueue(c)
	}
}

func (l *LifoScheduler) Pop() (c []byte) {
	if atomic.LoadInt64(l.closed) == 1 {
		return
	}

	item := l.q.Dequeue()
	if item == nil {
		return
	}

	return item.([]byte)
}

func (l *LifoScheduler) Close() {
	atomic.StoreInt64(l.closed, 1)
}
