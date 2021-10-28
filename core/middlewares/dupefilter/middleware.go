package dupefilter

import (
	"sync"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Middleware 去重中间件
type Middleware struct {
	mapping sync.Map
	config  *Config
}

type Config struct{}

func NewMiddleware(cfg *Config) *Middleware {
	return &Middleware{
		mapping: sync.Map{},
		config:  cfg,
	}
}

// Next implement core.Middleware.Next
func (d *Middleware) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) (err error) {
		u := c.GetRequest().URL.String()
		_, ok := d.mapping.LoadOrStore(u, struct{}{})
		if ok {
			return
		}
		err = handleFunc(c)
		if err != nil {
			d.mapping.Delete(u)
		}
		return
	}
}
