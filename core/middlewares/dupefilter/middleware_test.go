package dupefilter

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
	t.Run("未过滤", func(t *testing.T) {
		var flag bool
		err := filter.Next(func(c *core.Context) error {
			flag = true
			return nil
		})(c)
		assert.Equal(t, err, nil)
		assert.Equal(t, flag, true)
	})
	t.Run("过滤", func(t *testing.T) {
		var flag bool
		err := filter.Next(func(c *core.Context) error {
			flag = true
			return nil
		})(c)
		assert.Equal(t, err, nil)
		assert.Equal(t, flag, false)
	})
}
