package core

import (
	"net/http"
)

type RequestOptions struct {
	DontFilter *bool
	PreContext *Context
}

type RequestOptionsFunc func(options *RequestOptions)

// WithPreContext 上一次的Context
func WithPreContext(c *Context) RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.PreContext = c
	}
}

// DontFilter 不过滤
func DontFilter() RequestOptionsFunc {
	return func(options *RequestOptions) {
		dontFilter := true
		options.DontFilter = &dontFilter
	}
}

// Get 创建一个Get请求
func (e *Engine) Get(url string, handler HandlerFunc, options ...RequestOptionsFunc) (err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	e.Request(req, handler, options...)
	return
}

// Request 创建一个请求
func (e *Engine) Request(r *http.Request, handler HandlerFunc, options ...RequestOptionsFunc) {
	middlewares := append(e.middlewares, NewDownloader())

	middlewareHandler := e.applyMiddleware(func(c *Context) error {
		return nil
	}, middlewares...)

	ctx := NewContext(e, r, func(c *Context) error {
		if err := middlewareHandler(c); err != nil {
			return err
		}
		// 过滤等正常情况导致Response无数据
		if c.Response == nil {
			return nil
		}
		return handler(c)
	})

	opt := RequestOptions{}
	for _, optFunc := range options {
		optFunc(&opt)
	}

	if opt.PreContext != nil {
		ctx.SetDepth(opt.PreContext.GetDepth() + 1)
	}

	e.wg.Add(1)
	e.scheduler.Push(ctx)
}
