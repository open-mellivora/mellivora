package core

import (
	"testing"
)

func BenchmarkLifoScheduler(b *testing.B) {
	scheduler := NewLifoScheduler()
	cs := make([]*Context, b.N)
	for i := 0; i < b.N; i++ {
		cs[i] = NewContext(nil, nil, nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scheduler.Push(cs[i])
	}

	for i := 0; i < b.N; i++ {
		c := scheduler.Pop()
		if c != cs[i] {
			b.Errorf("c != ctxs[%d]", i)
		}
	}
	b.Run("空队列获取", func(b *testing.B) {
		if scheduler.Pop() != nil {
			b.Errorf("空队列有数据")
		}
	})

	b.Run("关闭后不能获取", func(b *testing.B) {
		scheduler.Close()
		if scheduler.Pop() != nil {
			b.Errorf("scheduler已关闭")
		}
	})
}
