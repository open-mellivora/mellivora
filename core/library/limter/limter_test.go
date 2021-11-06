package limter

import (
	"testing"
	"time"
)

func TestConcurrencyGroupLimiter_Wait(t *testing.T) {
	limiter := NewConcurrencyGroupLimiter(2)
	t.Run("未达到最大限制", func(t *testing.T) {
		limiter.Wait("a")
		limiter.Wait("a")
	})

	t.Run("达到最大限制", func(t *testing.T) {
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
	})
}

func TestConcurrencyGroupLimiter_Done(t *testing.T) {
	limiter := NewConcurrencyGroupLimiter(2)
	t.Run("在Wait之前Done", func(t *testing.T) {
		limiter.Done("a")
	})

	t.Run("Done后Wait", func(t *testing.T) {
		limiter.Wait("a")
		limiter.Wait("a")
		limiter.Done("a")
		limiter.Wait("a")
	})
}

func TestDelayGroupLimiter_Reset(t *testing.T) {
	t.Run("Reset未设置的Key", func(t *testing.T) {
		limiter := NewDelayGroupLimiter(time.Second)
		limiter.Reset("a")
	})

	t.Run("Reset存在的key", func(t *testing.T) {
		limiter := NewDelayGroupLimiter(time.Second)
		limiter.Wait("a")
		limiter.Reset("a")
		limiter.Wait("a")
	})
}

func TestDelayGroupLimiter_Wait(t *testing.T) {
	t.Run("延时0", func(t *testing.T) {
		limiter := NewDelayGroupLimiter(0)
		limiter.Wait("a")
		limiter.Wait("a")
		limiter.Wait("a")
	})

	t.Run("延时非0", func(t *testing.T) {
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
