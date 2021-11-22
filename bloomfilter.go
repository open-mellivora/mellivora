package mellivora

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilter struct {
	*bloom.BloomFilter
}

func (b *BloomFilter) Add(c *Context) {
	u := c.request.URL.String()
	b.BloomFilter.Add([]byte(u))
}

func NewBloomFilter() *BloomFilter {
	return &BloomFilter{
		BloomFilter: bloom.NewWithEstimates(1000000, 0.01),
	}
}

func (b *BloomFilter) Exist(c *Context) bool {
	u := c.request.URL.String()
	return b.BloomFilter.Test([]byte(u))
}
