package middlewares

import (
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Logging log access
type Logging struct {
	config LoggingConfig
}

// LoggingConfig defines the config for Logging middleware.
type LoggingConfig struct{}

// DefaultLoggingConfig is the default Logging middleware config.
var DefaultLoggingConfig struct{}

// NewLogging returns a Logging instance
func NewLogging() *Logging {
	return NewLoggingWithConfig(DefaultLoggingConfig)
}

// NewLoggingWithConfig returns a Logging middleware with config.
// See: `NewLogging()`.
func NewLoggingWithConfig(config LoggingConfig) *Logging {
	return &Logging{config: config}
}

// Next implement core.Logging.Next
func (m *Logging) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
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
