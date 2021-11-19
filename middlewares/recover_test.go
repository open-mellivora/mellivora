package middlewares

import (
	"net/http"
	"testing"

	"github.com/open-mellivora/mellivora"

	"github.com/go-playground/assert/v2"
)

func TestRecover(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	c := mellivora.NewContext(mellivora.NewEngine(32), req, nil)
	filter := NewRecover()
	t.Run("panic", func(t *testing.T) {
		err := filter.Next(func(c *mellivora.Context) error {
			panic("1")
		})(c)
		assert.NotEqual(t, err, nil)
	})
	t.Run("not panic", func(t *testing.T) {
		err := filter.Next(func(c *mellivora.Context) error {
			return nil
		})(c)
		assert.Equal(t, err, nil)
	})
}
