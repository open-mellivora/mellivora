package middlewares

import (
	"github.com/open-mellivora/mellivora"
)

// Retry defines a middleware which Retrys from panics anywhere in the chain
type Retry struct {
	RetryConfig
}

// RetryConfig defines the config for Retry middleware.
type RetryConfig struct {
	RetryTimes int
}

// NewRetryWithConfig returns a Retry middleware with config.
// See: `NewRetry()`.
func NewRetryWithConfig(config RetryConfig) *Retry {
	return &Retry{RetryConfig: config}
}

// DefaultRetryConfig is the default Retry middleware config.
var DefaultRetryConfig = RetryConfig{
	RetryTimes: 3,
}

// NewRetry returns a Retry instance
func NewRetry() *Retry {
	return NewRetryWithConfig(DefaultRetryConfig)
}

// Next implement mellivora.Middleware.Next
func (m *Retry) Next(handleFunc mellivora.MiddlewareFunc) mellivora.MiddlewareFunc {
	return func(c *mellivora.Context) (err error) {
		for i := 0; i < m.RetryTimes; i++ {
			err = handleFunc(c)
			if err != nil {
				continue
			} else {
				return
			}
		}
		return
	}
}
