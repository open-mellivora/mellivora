package mellivora

// Middleware defines a interface to process middleware.
type Middleware interface {
	Next(MiddlewareHandlerFunc) MiddlewareHandlerFunc
}

type middleware struct {
	next func(MiddlewareHandlerFunc) MiddlewareHandlerFunc
}

func (m *middleware) Next(handleFunc MiddlewareHandlerFunc) MiddlewareHandlerFunc {
	return m.next(handleFunc)
}

//NewMiddleware create a middleware instance
//nolint
func NewMiddleware(next func(MiddlewareHandlerFunc) MiddlewareHandlerFunc) *middleware {
	return &middleware{next: next}
}

var middlewares = make(map[string]Middleware)

func RegisterMiddleware(name string, m Middleware) {
	middlewares[name] = m
}

func GetMiddleware(name string) Middleware {
	return middlewares[name]
}
