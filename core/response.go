package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

type Response struct {
	*http.Response
	bodyBytes    []byte
	readBodyOnce sync.Once
}

func NewResponse(response *http.Response) *Response {
	return &Response{Response: response}
}

// Tokenizer get *html.Tokenizer from response.Body
func (resp *Response) Tokenizer() (tokenizer *html.Tokenizer, err error) {
	var bs []byte
	if bs, err = resp.Bytes(); err != nil {
		return nil, err
	}
	return html.NewTokenizer(bytes.NewBuffer(bs)), nil
}

// Bytes get []byte from response.Body
func (resp *Response) Bytes() (bodyBytes []byte, err error) {
	resp.readBodyOnce.Do(func() {
		if resp.bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		defer resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(resp.bodyBytes))
	})
	return resp.bodyBytes, err
}

// Bytes get []byte from response.Body
func (resp *Response) String() (str string, err error) {
	var bs []byte
	if bs, err = resp.Bytes(); err != nil {
		return
	}
	return string(bs), err
}

// JSON parses the resp.Body data and stores the result
func (resp *Response) JSON(i interface{}) error {
	bodyBytes, err := resp.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bodyBytes, i)
}
