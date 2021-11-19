package mellivora

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
	c.handler = handler.(HandleFunc)

	if c.request, err = http.ReadRequest(
		bufio.NewReader(bytes.NewBuffer(csz.RequestBytes))); err != nil {
		return
	}

	c.request.RequestURI = ""
	c.request.URL = csz.URL
	c.setter = csz.Setter
	return
}
