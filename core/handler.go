package core

// HandlerFunc defines a function to serve *Context.
type HandlerFunc func(c *Context) error

// Middleware defines a interface to process middleware.
type Middleware interface {
	Next(HandlerFunc) HandlerFunc
}

// Starter is the interface that wraps the basic Start method.
type Starter interface {
	Start(c *Engine)
}

// Closer is the interface that wraps the basic Close method.
type Closer interface {
	Close(c *Engine)
}

type middleware struct {
	next func(HandlerFunc) HandlerFunc
}

func (m *middleware) Next(handleFunc HandlerFunc) HandlerFunc {
	return m.next(handleFunc)
}

//NewMiddleware create a middleware instance
//nolint
func NewMiddleware(next func(HandlerFunc) HandlerFunc) *middleware {
	return &middleware{next: next}
}
