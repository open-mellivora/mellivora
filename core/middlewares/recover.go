package middlewares

import (
	"fmt"

	"github.com/pkg/errors"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Recover defines a middleware which recovers from panics anywhere in the chain
type Recover struct {
	RecoverConfig
}

// RecoverConfig defines the config for Recover middleware.
type RecoverConfig struct {
	RecoveryHandlerFunc HandlerFunc
}

// NewRecoverWithConfig returns a Recover middleware with config.
// See: `NewRecover()`.
func NewRecoverWithConfig(config RecoverConfig) *Recover {
	return &Recover{RecoverConfig: config}
}

// DefaultRecoverConfig is the default Recover middleware config.
var DefaultRecoverConfig = RecoverConfig{
	RecoveryHandlerFunc: func(p interface{}) (err error) {
		return errors.WithStack(fmt.Errorf("%v", p))
	},
}

// NewRecover returns a Recover instance
func NewRecover() *Recover {
	return NewRecoverWithConfig(DefaultRecoverConfig)
}

// HandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type HandlerFunc func(p interface{}) (err error)

// Next implement core.Middleware.Next
func (m *Recover) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) (err error) {
		panicked := true
		defer func() {
			if r := recover(); r != nil || panicked {
				err = m.RecoverConfig.RecoveryHandlerFunc(r)
				c.Engine().Logger().Error("Recovery error: %+v", err)
			}
		}()

		err = handleFunc(c)
		panicked = false
		return
	}
}
