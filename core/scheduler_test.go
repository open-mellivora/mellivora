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
			b.Errorf("c != ctxs[%d],c:%v,cs[%d]:%v", i, c, i, cs[i])
		}
	}
	b.StopTimer()
}

func TestLifoScheduler(t *testing.T) {
	scheduler := NewLifoScheduler()
	t.Run("空队列获取", func(t *testing.T) {
		if scheduler.Pop() != nil {
			t.Errorf("空队列有数据")
		}
	})
	scheduler.Push(new(Context))
	t.Run("关闭后不能获取", func(t *testing.T) {
		scheduler.Close()
		if scheduler.Pop() != nil {
			t.Errorf("scheduler已关闭")
		}
	})
}
