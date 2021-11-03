package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"sync"

	"github.com/valyala/bytebufferpool"
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
		request    *http.Request
		core       *Engine
		handler    HandlerFunc
		httpClient *http.Client
		setter
	}
)

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
		if err = c.core.request(c, req, handler, options...); err != nil {
			return
		}
	}

	return nil
}

type ContextSerializable struct {
	HandlerName  string
	RequestBytes []byte
	URL          *url.URL
	Setter       setter
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

func (cs *ContextSerializer) Marshal(c *Context) (bs []byte, err error) {
	handlerName := reflect.TypeOf(c.handler).String()
	cs.handlers.LoadOrStore(handlerName, c.handler)
	csz := &ContextSerializable{
		HandlerName: handlerName,
		URL:         c.GetRequest().URL,
		Setter:      c.setter,
	}
	if csz.RequestBytes, err = httputil.DumpRequest(c.GetRequest(), true); err != nil {
		return
	}
	return json.Marshal(csz)
}

func (cs *ContextSerializer) Unmarshal(bs []byte) (c *Context, err error) {
	var csz ContextSerializable
	if err = json.Unmarshal(bs, &csz); err != nil {
		return
	}
	c = new(Context)
	var ok bool
	var handler interface{}
	if handler, ok = cs.handlers.Load(csz.HandlerName); !ok {
		return
	}
	c.handler = handler.(HandlerFunc)

	if c.request, err = http.ReadRequest(
		bufio.NewReader(bytes.NewBuffer(csz.RequestBytes))); err != nil {
		return
	}
	c.request.RequestURI = ""
	c.request.URL = csz.URL
	c.setter = csz.Setter
	return
}
