package core

type HandleFunc func(c *Context) error

type Middleware interface {
	Next(HandleFunc) HandleFunc
}

type Starter interface {
	// Start call if set when core start
	Start(c *Engine)
}

type Closer interface {
	// Close call if set when core close
	Close(c *Engine)
}

type middleware struct {
	next func(HandleFunc) HandleFunc
}

func (m *middleware) Next(handleFunc HandleFunc) HandleFunc {
	return m.next(handleFunc)
}

//nolint
func NewMiddleware(next func(HandleFunc) HandleFunc) *middleware {
	return &middleware{next: next}
}
