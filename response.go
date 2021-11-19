package mellivora

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"
)

type Response struct {
	*http.Response
}

func NewResponse(response *http.Response) *Response {
	return &Response{Response: response}
}

// Tokenizer get *html.Tokenizer from response.Body
func (resp *Response) Tokenizer() (tokenizer *html.Tokenizer) {
	return html.NewTokenizer(resp.Body)
}

// Bytes get []byte from response.Body
func (resp *Response) Bytes() (bodyBytes []byte, err error) {
	return ioutil.ReadAll(resp.Body)
}

// Bytes get []byte from response.Body
func (resp *Response) String() (str string, err error) {
	var bs []byte
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return string(bs), err
}

// JSON parses the resp.Body data and stores the result
func (resp *Response) JSON(i interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(i)
}
