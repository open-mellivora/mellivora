package core

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRequestOptions_DontFilter(t *testing.T) {
	opt := NewRequestOptions()
	opts := []RequestOptionsFunc{DontFilter()}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	assert.Equal(t, opt.GetDontFilter(), true)
}

func TestRequestOptions_WithDepth(t *testing.T) {
	opt := NewRequestOptions()
	opts := []RequestOptionsFunc{withDepth(3)}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	assert.Equal(t, opt.GetDepth(), int64(3))
}

func TestRequestOptions_WithValue(t *testing.T) {
	opt := NewRequestOptions()
	opts := []RequestOptionsFunc{WithValue("a", "b")}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	assert.Equal(t, opt.Value("a"), "b")
}
