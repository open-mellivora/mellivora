package middlewares

import (
	"github.com/open-mellivora/mellivora"
	"go.uber.org/zap"
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

// Next implement mellivora.Next
func (m *Logging) Next(handleFunc mellivora.MiddlewareFunc) mellivora.MiddlewareFunc {
	return func(c *mellivora.Context) (err error) {
		err = handleFunc(c)
		logger := c.Engine().Logger().With(
			zap.String("url", c.GetRequest().URL.String()),
			zap.Int64("depth", c.GetDepth()))

		if err != nil {
			logger.Warn("logging error", zap.Error(err))
			return
		}

		statusCode := c.Response.StatusCode
		logger.Info("logging", zap.Int("statusCode", statusCode))
		return err
	}
}
