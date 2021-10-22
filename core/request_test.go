package core

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRequestOptions(t *testing.T) {
	opt := RequestOptions{}
	c := NewContext(nil, nil, nil)
	opts := []RequestOptionsFunc{WithPreContext(c), DontFilter()}
	for _, optFunc := range opts {
		optFunc(&opt)
	}
	assert.Equal(t, opt.DontFilter, true)
	assert.Equal(t, opt.PreContext, c)
}
