package core

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"sync"

	"github.com/valyala/bytebufferpool"
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
		Setter
	}
)

type Setter struct {
	m map[interface{}]interface{}
}

func newSetter() Setter {
	return Setter{m: make(map[interface{}]interface{})}
}

func (c *Setter) Set(k, v interface{}) {
	c.m[k] = v
}

func (c *Setter) Value(k interface{}) interface{} {
	return c.m[k]
}

// NewContext returns a Context instance.
func NewContext(core *Engine, request *http.Request, handler HandlerFunc) *Context {
	return &Context{
		request: request,
		core:    core,
		handler: handler,
		Setter:  newSetter(),
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
	c.Setter.Set(depthKey, depth)
}

// GetDepth returns `depth`.
func (c *Context) GetDepth() int64 {
	if c == nil {
		return 0
	}
	value := c.Setter.Value(depthKey)
	if value == nil {
		return 0
	}
	return value.(int64)
}

type ContextSerializable struct {
	HandlerName  string
	RequestBytes []byte
	URL          *url.URL
	Setter       map[interface{}]interface{}
}

type ContextSerializer struct {
	pool     *bytebufferpool.Pool
	handlers sync.Map
}

func NewContextSerializer() *ContextSerializer {
	return &ContextSerializer{
		pool: new(bytebufferpool.Pool),
	}
}

func (cs *ContextSerializer) Marshal(c *Context) (csz *ContextSerializable, err error) {
	handlerName := reflect.TypeOf(c.handler).String()
	cs.handlers.LoadOrStore(handlerName, c.handler)
	csz = &ContextSerializable{
		HandlerName: handlerName,
		URL:         c.GetRequest().URL,
		Setter:      c.Setter.m,
	}

	if csz.RequestBytes, err = httputil.DumpRequest(c.GetRequest(), true); err != nil {
		return
	}
	return
}

func (cs *ContextSerializer) Unmarshal(csz *ContextSerializable) (c *Context, err error) {
	c = new(Context)
	var ok bool
	var handler interface{}
	if handler, ok = cs.handlers.Load(csz.HandlerName); !ok {
		return
	}
	c.handler = handler.(HandlerFunc)
	c.m = csz.Setter

	if c.request, err = http.ReadRequest(bufio.NewReader(bytes.NewBuffer(csz.RequestBytes))); err != nil {
		return
	}
	c.request.RequestURI = ""
	c.request.URL = csz.URL
	return
}
