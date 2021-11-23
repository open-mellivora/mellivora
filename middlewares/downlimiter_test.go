package middlewares

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"github.com/open-mellivora/mellivora"
)

func TestDownLimiter(t *testing.T) {
	limiter := NewDownLimiterWithConfig(DownLimiterConfig{Timeout: time.Second})
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	c := mellivora.NewContext(nil, req, nil)
	err := limiter.Next(func(c *mellivora.Context) error {
		_, ok := c.GetRequest().Context().Deadline()
		assert.Equal(t, ok, true)
		return nil
	})(c)
	assert.Equal(t, err, nil)
}
