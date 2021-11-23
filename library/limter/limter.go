package limter

import (
	"sync"
	"time"
)

type ConcurrencyLimiter struct {
	c chan struct{}
}

// NewConcurrencyLimiter returns a ConcurrencyLimiter instance
// up to the given maximum capacity
func NewConcurrencyLimiter(capacity uint64) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{c: make(chan struct{}, capacity)}
}

// Wait blocking request tokens until they are available.
func (l *ConcurrencyLimiter) Wait() {
	l.c <- struct{}{}
}

// Done release a token
func (l *ConcurrencyLimiter) Done() {
	<-l.c
}

type ConcurrencyGroupLimiter struct {
	mapping  sync.Map
	capacity uint64
}

// NewConcurrencyGroupLimiter returns a ConcurrencyGroupLimiter instance
func NewConcurrencyGroupLimiter(capacity uint64) *ConcurrencyGroupLimiter {
	return &ConcurrencyGroupLimiter{
		mapping:  sync.Map{},
		capacity: capacity,
	}
}

// Wait blocking request tokens until they are available.
func (l *ConcurrencyGroupLimiter) Wait(key string) {
	limiter := NewConcurrencyLimiter(l.capacity)
	v, ok := l.mapping.LoadOrStore(key, limiter)
	if ok {
		v.(*ConcurrencyLimiter).Wait()
	} else {
		limiter.Wait()
	}
}

// Done release a token
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

// Wait blocking request tokens until they are available.
func (l *DelayGroupLimiter) Wait(key string) {
	if l.delay <= 0 {
		return
	}
	limiter := time.NewTimer(l.delay)
	v, ok := l.mapping.LoadOrStore(key, limiter)
	if ok {
		<-v.(*time.Timer).C
	}
}

// Reset changes the timer to expire after duration delay.
func (l *DelayGroupLimiter) Reset(key string) {
	if v, ok := l.mapping.Load(key); ok {
		v.(*time.Timer).Reset(l.delay)
	} else {
		return
	}
}
