package middlewares

import (
	"io"
	"io/ioutil"

	"github.com/open-mellivora/mellivora"

	"golang.org/x/net/html/charset"
)

// Decoder for decode response body to utf8
type Decoder struct {
	config DecoderConfig
}

// DecoderConfig defines the config for Decoder middleware.
type DecoderConfig struct{}

// DefaultDecoderConfig is the default Decoder middleware config.
var DefaultDecoderConfig struct{}

// NewDecoder returns a Decoder instance
func NewDecoder() *Decoder {
	return NewDecoderWithConfig(DefaultDecoderConfig)
}

// NewDecoderWithConfig returns a Decoder middleware with config.
// See: `NewDecoder()`.
func NewDecoderWithConfig(config DecoderConfig) *Decoder {
	return &Decoder{config: config}
}

// Next implement mellivora.Middleware.Next
func (p *Decoder) Next(handleFunc mellivora.MiddlewareHandlerFunc) mellivora.MiddlewareHandlerFunc {
	return func(c *mellivora.Context) (err error) {
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

func init() {
	mellivora.RegisterMiddleware("a", &Decoder{})
}
