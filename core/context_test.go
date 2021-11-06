package core

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestContext_Set(t *testing.T) {
	e := NewEngine(32)
	req := new(http.Request)
	handler := func(c *Context) error {
		return nil
	}

	c := NewContext(e, req, handler)
	assert.Equal(t, req, c.GetRequest())
	assert.Equal(t, e, c.Engine())

	req = new(http.Request)
	c.SetRequest(req)
	assert.Equal(t, req, c.GetRequest())

	resp := NewResponse(nil)
	c.SetResponse(resp)
	assert.Equal(t, resp, c.Response)

	client := &http.Client{}
	c.SetHTTPClient(client)
	assert.Equal(t, c.httpClient, client)
}
