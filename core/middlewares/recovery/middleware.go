package recovery

import (
	"fmt"

	"github.com/pkg/errors"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

type Middleware struct {
	recoveryHandlerFunc HandlerFunc
}

func NewMiddleware(f HandlerFunc) *Middleware {
	if f == nil {
		f = func(p interface{}) (err error) {
			return errors.WithStack(fmt.Errorf("%v", p))
		}
	}
	return &Middleware{recoveryHandlerFunc: f}
}

// HandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type HandlerFunc func(p interface{}) (err error)

func (m *Middleware) Next(handleFunc core.HandleFunc) core.HandleFunc {
	return func(c *core.Context) (err error) {
		panicked := true
		defer func() {
			if r := recover(); r != nil || panicked {
				err = m.recoveryHandlerFunc(r)
				c.Core().Logger().Error("Recovery error: %+v", err)
			}
		}()

		err = handleFunc(c)
		panicked = false
		return
	}
}
