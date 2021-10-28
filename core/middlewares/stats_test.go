package middlewares

import (
	"net/http"
	"sort"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/pkg/errors"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

func TestStatsCollector(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	e := core.NewEngine()
	c := core.NewContext(e, req, nil)
	filter := NewStatsCollector()
	t.Run("记录信息", func(t *testing.T) {
		err := filter.Next(func(c *core.Context) error {
			c.SetResponse(core.NewResponse(new(http.Response)))
			return nil
		})(c)
		assert.Equal(t, err, nil)
	})
	t.Run("记录错误信息", func(t *testing.T) {
		err := filter.Next(func(c *core.Context) error {
			c.SetResponse(core.NewResponse(new(http.Response)))
			return errors.New("x")
		})(c)
		assert.NotEqual(t, err, nil)
	})
	t.Run("关闭", func(t *testing.T) {
		filter.Close(e)
	})
}

func TestSortKVs(t *testing.T) {
	kvs := sortKVS{kv{key: "l1/la"}, kv{key: "l1"}, kv{key: "l2"}, kv{key: "l1/lb"}}
	sort.Sort(kvs)
	assert.Equal(t, kvs[0].key, "l1")
	assert.Equal(t, kvs[1].key, "l1/la")
	assert.Equal(t, kvs[2].key, "l1/lb")
	assert.Equal(t, kvs[3].key, "l2")
}
