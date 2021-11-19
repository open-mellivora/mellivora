package mellivora

type RequestOptions struct {
	setter
}

func NewRequestOptions() *RequestOptions {
	return &RequestOptions{setter: newSetter()}
}

type RequestOptionsFunc func(options *RequestOptions)

// withDepth returns a RequestOptionsFunc which sets the depth
func withDepth(depth int64) RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.setter.SetDepth(depth)
	}
}

// DontFilter returns a RequestOptionsFunc which sets the dontFilter
func DontFilter() RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.setter.SetDontFilter(true)
	}
}

// WithValue returns a RequestOptionsFunc which sets k,v in setter
// The provided value must be serializable
func WithValue(k string, v interface{}) RequestOptionsFunc {
	return func(options *RequestOptions) {
		options.setter.Set(k, v)
	}
}
