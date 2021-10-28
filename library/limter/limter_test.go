package limter

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

func TestDelayGroupLimiter_Wait(t *testing.T) {
	t.Run("可重入", func(t *testing.T) {
		limiter := NewDelayGroupLimiter(0)
		limiter.Wait("a")
		limiter.Wait("a")
		limiter.Wait("a")
	})

	t.Run("reset有效", func(t *testing.T) {
		limiter := NewDelayGroupLimiter(time.Second)
		limiter.Wait("a")
		limiter.Reset("a")
		limiter.Wait("a")
	})
	t.Run("延时生效", func(t *testing.T) {
		limiter := NewDelayGroupLimiter(time.Second)
		limiter.Wait("a")
		c := make(chan struct{})
		go func() {
			limiter.Wait("a")
			c <- struct{}{}
		}()
		select {
		case <-c:
			t.Error("延时失败")
		case <-time.After(300 * time.Microsecond):
		}
	})
}
