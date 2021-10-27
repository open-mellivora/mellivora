package downlimiter

import (
	"sync"
	"time"
)

type ConcurrencyLimiter struct {
	c chan struct{}
}

func NewConcurrencyLimiter(size int64) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{c: make(chan struct{}, size)}
}

// Wait 阻塞等待
func (l *ConcurrencyLimiter) Wait() {
	l.c <- struct{}{}
}

// Done 完成
func (l *ConcurrencyLimiter) Done() {
	<-l.c
}

type ConcurrencyGroupLimiter struct {
	mapping sync.Map
	size    int64
}

func NewConcurrencyGroupLimiter(size int64) *ConcurrencyGroupLimiter {
	return &ConcurrencyGroupLimiter{
		mapping: sync.Map{},
		size:    size,
	}
}

// Wait 阻塞等待
func (l *ConcurrencyGroupLimiter) Wait(key string) {
	limiter := NewConcurrencyLimiter(l.size)
	v, ok := l.mapping.LoadOrStore(key, limiter)
	if ok {
		v.(*ConcurrencyLimiter).Wait()
	} else {
		limiter.Wait()
	}
}

// Done 完成
func (l *ConcurrencyGroupLimiter) Done(key string) {
	if v, ok := l.mapping.Load(key); ok {
		v.(*ConcurrencyLimiter).Done()
	} else {
		return
	}
}

type DelayGroupLimiter struct {
	mapping sync.Map
	delay   time.Duration
}

func NewDelayGroupLimiter(delay time.Duration) *DelayGroupLimiter {
	return &DelayGroupLimiter{delay: delay}
}

// Wait 阻塞等待
func (l *DelayGroupLimiter) Wait(key string) {
	if l.delay == 0 {
		return
	}
	limiter := time.NewTimer(l.delay)
	v, ok := l.mapping.LoadOrStore(key, limiter)
	if ok {
		<-v.(*time.Timer).C
	}
}

// Reset 重置定时器
func (l *DelayGroupLimiter) Reset(key string) {
	if v, ok := l.mapping.Load(key); ok {
		v.(*time.Timer).Reset(l.delay)
	} else {
		return
	}
}
