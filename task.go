package mellivora

import "net/http"

type Task interface {
	Requests() []*http.Request
	Items() []interface{}
	RequestOptions() []RequestOptionsFunc
	Handler() HandleFunc
}

func NewItems(items ...interface{}) Task {
	return task{
		items: items,
	}
}

func Gets(urls []string, handler HandleFunc, options ...RequestOptionsFunc) (t Task, err error) {
	requests := make([]*http.Request, len(urls))
	for i := 0; i < len(urls); i++ {
		if requests[i], err = http.NewRequest(http.MethodGet, urls[i], nil); err != nil {
			return
		}
	}
	return task{
		requests:       requests,
		requestOptions: options,
		handler:        handler,
	}, nil
}

func Request(request *http.Request, handler HandleFunc, options ...RequestOptionsFunc) Task {
	return task{
		requests:       []*http.Request{request},
		requestOptions: options,
		handler:        handler,
	}
}

func MustGet(url string, handler HandleFunc, options ...RequestOptionsFunc) Task {
	if task, err := Get(url, handler, options...); err != nil {
		panic(err)
	} else {
		return task
	}
}

func Get(url string, handler HandleFunc, options ...RequestOptionsFunc) (Task, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return Request(req, handler, options...), nil
}

type task struct {
	requests       []*http.Request
	items          []interface{}
	requestOptions []RequestOptionsFunc
	handler        HandleFunc
}

func (t task) Requests() []*http.Request {
	return t.requests
}

func (t task) Items() []interface{} {
	return t.items
}

func (t task) RequestOptions() []RequestOptionsFunc {
	return t.requestOptions
}

func (t task) Handler() HandleFunc {
	return t.handler
}
