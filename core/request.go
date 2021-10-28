package core

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
