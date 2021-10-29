package core

import (
	"net/http"
)

const depthKey depth = iota

type (
	depth int64
	// Context represents the context of the current HTTP request. It holds request and
	// response objects, path, path parameters, data and registered handler.
	Context struct {
		*Response
		request    *http.Request
		core       *Engine
		handler    HandlerFunc
		httpClient *http.Client
		setter
	}
)

type setter struct {
	m map[interface{}]interface{}
}

func newSetter() setter {
	return setter{m: make(map[interface{}]interface{})}
}

func (c *setter) Set(k, v interface{}) {
	c.m[k] = v
}

func (c *setter) Value(k interface{}) interface{} {
	return c.m[k]
}

// NewContext returns a Context instance.
func NewContext(core *Engine, request *http.Request, handler HandlerFunc) *Context {
	return &Context{
		request: request,
		core:    core,
		handler: handler,
		setter:  newSetter(),
	}
}

// SetHTTPClient sets `*http.Client`.
func (c *Context) SetHTTPClient(client *http.Client) {
	c.httpClient = client
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

// SetDepth sets `depth`.
func (c *Context) SetDepth(depth int64) {
	if c == nil {
		return
	}
	c.setter.Set(depthKey, depth)
}

// GetDepth returns `depth`.
func (c *Context) GetDepth() int64 {
	if c == nil {
		return 0
	}
	value := c.setter.Value(depthKey)
	if value == nil {
		return 0
	}
	return value.(int64)
}
