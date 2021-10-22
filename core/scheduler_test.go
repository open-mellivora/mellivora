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
		c := scheduler.BlockPop()
		if c != cs[i] {
			b.Errorf("c != ctxs[%d]", i)
		}
	}
	scheduler.Close()
	if scheduler.BlockPop() != nil {
		b.Errorf("scheduler已关闭")
	}
}
