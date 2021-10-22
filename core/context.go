package core

import (
	"context"
	"net/http"
)

const depthKey depth = iota

type (
	depth int64
	// Context 请求处理上下文
	Context struct {
		*Response
		request    *http.Request
		core       *Engine
		handler    HandleFunc
		httpClient *http.Client
		setter
	}
)

type setter struct {
	ctx context.Context
}

func newSetter() setter {
	return setter{ctx: context.TODO()}
}

func (c *setter) Set(k, v interface{}) {
	c.ctx = context.WithValue(c.ctx, k, v)
}

func (c *setter) Value(k interface{}) interface{} {
	return c.ctx.Value(k)
}

// NewContext create a Context
func NewContext(core *Engine, request *http.Request, handler HandleFunc) *Context {
	return &Context{
		request: request,
		core:    core,
		handler: handler,
		setter:  newSetter(),
	}
}

// SetHTTPClient 设置http.Client
func (c *Context) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *Context) SetResponse(response *Response) {
	c.Response = response
}

func (c *Context) Core() *Engine {
	return c.core
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) SetRequest(req *http.Request) {
	c.request = req
}

func (c *Context) SetDepth(depth int64) {
	if c == nil {
		return
	}
	c.setter.Set(depthKey, depth)
}

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
