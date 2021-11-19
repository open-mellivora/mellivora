package middlewares

import (
	"net/http"
	"testing"

	"github.com/open-mellivora/mellivora"

	"github.com/go-playground/assert/v2"
	"github.com/pkg/errors"
)

func TestLogging(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	e := mellivora.NewEngine(32)
	c := mellivora.NewContext(e, req, nil)
	filter := NewLogging()
	t.Run("记录信息", func(t *testing.T) {
		err := filter.Next(func(c *mellivora.Context) error {
			c.SetResponse(mellivora.NewResponse(new(http.Response)))
			return nil
		})(c)
		assert.Equal(t, err, nil)
	})
	t.Run("记录错误信息", func(t *testing.T) {
		err := filter.Next(func(c *mellivora.Context) error {
			c.SetResponse(mellivora.NewResponse(new(http.Response)))
			return errors.New("x")
		})(c)
		assert.NotEqual(t, err, nil)
	})
}
