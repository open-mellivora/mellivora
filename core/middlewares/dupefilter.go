package middlewares

import (
	"sync"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// DupeFilter filter same url
type DupeFilter struct {
	config DupeFilterConfig
	m      sync.Map
}

// DupeFilterConfig defines the config for DupeFilter middleware.
type DupeFilterConfig struct{}

// DefaultDupeFilterConfig is the default DupeFilter middleware config.
var DefaultDupeFilterConfig struct{}

// NewDupeFilter returns a DupeFilter instance
func NewDupeFilter() *DupeFilter {
	return NewDupeFilterWithConfig(DefaultDupeFilterConfig)
}

// NewDupeFilterWithConfig returns a DupeFilter middleware with config.
// See: `NewDupeFilter()`.
func NewDupeFilterWithConfig(config DupeFilterConfig) *DupeFilter {
	return &DupeFilter{config: config}
}

// Next implement core.Middleware.Next
func (d *DupeFilter) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) (err error) {
		u := c.GetRequest().URL.String()
		_, ok := d.m.LoadOrStore(u, struct{}{})
		if ok {
			return
		}
		err = handleFunc(c)
		if err != nil {
			d.m.Delete(u)
		}
		return
	}
}
