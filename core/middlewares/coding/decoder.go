package coding

import (
	"io"
	"io/ioutil"

	"golang.org/x/net/html/charset"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Decoder 自动转response的编码为utf8
type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

// Next implement core.Middleware.Next
func (p *Decoder) Next(handleFunc core.HandleFunc) core.HandleFunc {
	return func(c *core.Context) (err error) {
		if err = handleFunc(c); err != nil {
			return
		}

		if c.Response == nil {
			return
		}

		var newReader io.Reader
		if newReader, err = charset.NewReader(
			c.Response.Body,
			c.Response.Header.Get("Context-Type"),
		); err != nil {
			return err
		}
		c.Response.Body = ioutil.NopCloser(newReader)
		return
	}
}
