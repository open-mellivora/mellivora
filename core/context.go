package core

import (
	"net/http"
)

const (
	depthKey      depth      = iota
	dontFilterKey dontFilter = iota
)

type (
	depth      int64
	dontFilter int64
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

// SetDontFilter sets `depth`.
func (c *setter) SetDontFilter(dontFilter bool) {
	c.Set(dontFilterKey, dontFilter)
}

// GetDontFilter returns `depth`.
func (c *setter) GetDontFilter() bool {
	if c == nil {
		return false
	}
	value := c.Value(dontFilterKey)
	if value == nil {
		return false
	}
	return value.(bool)
}

// SetDepth sets `depth`.
func (c *setter) SetDepth(depth int64) {
	c.Set(depthKey, depth)
}

// GetDepth returns `depth`.
func (c *setter) GetDepth() int64 {
	if c == nil {
		return 0
	}
	value := c.Value(depthKey)
	if value == nil {
		return 0
	}
	return value.(int64)
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

// Get create tasks from url
func (c *Context) Get(url string, handler HandlerFunc,
	options ...RequestOptionsFunc) (err error) {

	return c.Gets([]string{url}, handler, options...)
}

// Gets create tasks from urls
func (c *Context) Gets(urls []string, handler HandlerFunc,
	options ...RequestOptionsFunc) (err error) {

	reqs := make([]*http.Request, len(urls))
	for i := 0; i < len(urls); i++ {
		if reqs[i], err = http.NewRequest(http.MethodGet, urls[i], nil); err != nil {
			return
		}
	}

	return c.Requests(reqs, handler, options...)
}

// Request create task from req
func (c *Context) Request(req *http.Request, handler HandlerFunc,
	options ...RequestOptionsFunc) (err error) {

	return c.Requests([]*http.Request{req}, handler, options...)
}

// Requests create tasks from reqs
func (c *Context) Requests(reqs []*http.Request, handler HandlerFunc,
	options ...RequestOptionsFunc) (err error) {

	if c.request != nil {
		options = append(options, withDepth(c.GetDepth()+1))
	}

	for i := 0; i < len(reqs); i++ {
		req := reqs[i]
		c.core.request(c, req, handler, options...)
	}

	return nil
}
