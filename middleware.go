package mellivora

// Middleware defines a interface to process middleware.
type Middleware interface {
	Next(MiddlewareFunc) MiddlewareFunc
}

type middleware struct {
	next func(MiddlewareFunc) MiddlewareFunc
}

func (m *middleware) Next(handleFunc MiddlewareFunc) MiddlewareFunc {
	return m.next(handleFunc)
}

//NewMiddleware create a middleware instance
//nolint
func NewMiddleware(next func(MiddlewareFunc) MiddlewareFunc) *middleware {
	return &middleware{next: next}
}

var middlewares = make(map[string]Middleware)

func RegisterMiddleware(name string, m Middleware) {
	middlewares[name] = m
}

func GetMiddleware(name string) Middleware {
	return middlewares[name]
}
