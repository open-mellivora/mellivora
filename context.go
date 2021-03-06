package mellivora

import (
	"net/http"
)

const (
	depthKey      string = "_depth"
	dontFilterKey string = "_dontFilter"
)

type (
	// Context represents the context of the current HTTP request. It holds request and
	// response objects, path, path parameters, data and registered handler.
	Context struct {
		*Response
		request      *http.Request
		core         *Engine
		handler      HandleFunc
		roundTripper http.RoundTripper
		setter
	}
)

// NewContext returns a Context instance.
func NewContext(core *Engine, request *http.Request, handler HandleFunc) *Context {
	return &Context{
		request: request,
		core:    core,
		handler: handler,
		setter:  newSetter(),
	}
}

// SetRoundTripper sets `http.RoundTripper`.
func (c *Context) SetRoundTripper(rt http.RoundTripper) {
	c.roundTripper = rt
}

// SetResponse sets `*Response`.
func (c *Context) SetResponse(response *Response) {
	c.Response = response
}

// Engine returns the `Engine` instance.
func (c *Context) Engine() *Engine {
	return c.core
}

// GetRequest returns `*http.Request`.
func (c *Context) GetRequest() *http.Request {
	return c.request
}

// SetRequest sets `*http.Request`.
func (c *Context) SetRequest(req *http.Request) {
	c.request = req
}
