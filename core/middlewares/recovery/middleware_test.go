package recovery

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

func TestMiddleware_Next(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	c := core.NewContext(core.NewEngine(), req, nil)
	filter := NewMiddleware(nil)
	t.Run("panic", func(t *testing.T) {
		err := filter.Next(func(c *core.Context) error {
			panic("1")
		})(c)
		assert.NotEqual(t, err, nil)
	})
	t.Run("not panic", func(t *testing.T) {
		err := filter.Next(func(c *core.Context) error {
			return nil
		})(c)
		assert.Equal(t, err, nil)
	})
}
