package core

type RequestOptions struct {
	setter
}

func NewRequestOptions() *RequestOptions {
	return &RequestOptions{setter: newSetter()}
}

type RequestOptionsFunc func(options *RequestOptions)

// withDepth 设置depth
func withDepth(depth int64) RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.setter.SetDepth(depth)
	}
}

// DontFilter 不过滤
func DontFilter() RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.setter.SetDontFilter(true)
	}
}

func WithValue(k string, v interface{}) RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.setter.Set(k, v)
	}
}
