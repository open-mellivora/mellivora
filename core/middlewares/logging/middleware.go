package logging

import (
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Middleware 日志中间件
type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

// Next implement core.Middleware.Next
func (m *Middleware) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) (err error) {
		err = handleFunc(c)
		if err != nil {
			c.Engine().Logger().Error("[depth]:%v [url]:%v [error]:%v",
				c.GetDepth(), c.GetRequest().URL.String(), err)
			return
		}
		statusCode := c.Response.StatusCode
		c.Engine().Logger().Debug("[depth]:%v [url]:%v [status]:%v ",
			c.GetDepth(), c.GetRequest().URL.String(), statusCode)
		return err
	}
}
