package mellivora

import (
	"bufio"
	"bytes"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestResponse_String(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	resp, _ := http.ReadResponse(bufio.NewReader(bytes.NewBufferString(`HTTP/1.1 206 Partial Content
Connection: close
Content-Type: multipart/byteranges; boundary=18a75608c8f47cef

{"a":"1"}`)), req)
	c := NewContext(nil, req, nil)
	c.SetResponse(NewResponse(resp))

	str, err := c.String()
	assert.Equal(t, err, nil)
	assert.Equal(t, str, "{\"a\":\"1\"}")
}

func TestResponse_Bytes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	resp, _ := http.ReadResponse(bufio.NewReader(bytes.NewBufferString(`HTTP/1.1 206 Partial Content
Connection: close
Content-Type: multipart/byteranges; boundary=18a75608c8f47cef

{"a":"1"}`)), req)
	c := NewContext(nil, req, nil)
	c.SetResponse(NewResponse(resp))

	bs, err := c.Bytes()
	assert.Equal(t, err, nil)
	assert.Equal(t, bs, []byte("{\"a\":\"1\"}"))
}

func TestResponse_JSON(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	resp, _ := http.ReadResponse(bufio.NewReader(bytes.NewBufferString(`HTTP/1.1 206 Partial Content
Connection: close
Content-Type: multipart/byteranges; boundary=18a75608c8f47cef

{"a":"1"}`)), req)
	c := NewContext(nil, req, nil)
	c.SetResponse(NewResponse(resp))

	mapping := make(map[string]string)
	err := c.JSON(&mapping)
	assert.Equal(t, err, nil)
	assert.Equal(t, mapping["a"], "1")
}
