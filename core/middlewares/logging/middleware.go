package logging

import (
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Middleware 日志中间件
type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Next(handleFunc core.HandleFunc) core.HandleFunc {
	return func(c *core.Context) (err error) {
		err = handleFunc(c)
		if err != nil {
			c.Core().Logger().Error("[depth]:%v [url]:%v [error]:%v",
				c.GetDepth(), c.GetRequest().URL.String(), err)
			return
		}
		statusCode := c.Response.StatusCode
		c.Core().Logger().Debug("[depth]:%v [url]:%v [status]:%v ",
			c.GetDepth(), c.GetRequest().URL.String(), statusCode)
		return err
	}
}
