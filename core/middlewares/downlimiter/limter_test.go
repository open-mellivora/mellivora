package downlimiter

import (
	"testing"
	"time"
)

func TestConcurrencyGroupLimiter_Wait(t *testing.T) {
	limiter := NewConcurrencyGroupLimiter(2)
	limiter.Wait("a")
	limiter.Wait("a")
	c := make(chan struct{})
	go func() {
		limiter.Wait("a")
		c <- struct{}{}
	}()
	select {
	case <-c:
		t.Errorf("限制错误")
	case <-time.After(20 * time.Millisecond):
	}
	limiter.Done("a")
	limiter.Done("a")
	limiter.Wait("a")
}
