package middlewares

import (
	"encoding/base64"

	"github.com/open-mellivora/mellivora"
)

// BasicAuth for decode response body to utf8
type BasicAuth struct {
	config BasicAuthConfig
}

// BasicAuthConfig defines the config for BasicAuth middleware.
type BasicAuthConfig struct {
	UserName string
	Password string
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// DefaultBasicAuthConfig is the default BasicAuth middleware config.
var DefaultBasicAuthConfig = BasicAuthConfig{}

// NewBasicAuth returns a BasicAuth instance
func NewBasicAuth() *BasicAuth {
	return NewBasicAuthWithConfig(DefaultBasicAuthConfig)
}

// NewBasicAuthWithConfig returns a BasicAuth middleware with config.
// See: `NewBasicAuth()`.
func NewBasicAuthWithConfig(config BasicAuthConfig) *BasicAuth {
	return &BasicAuth{config: config}
}

// Next implement mellivora.Middleware.Next
func (p *BasicAuth) Next(handleFunc mellivora.MiddlewareFunc) mellivora.MiddlewareFunc {
	return func(c *mellivora.Context) (err error) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.Request.Header.Set("Authorization", "Basic "+basicAuth(p.config.UserName, p.config.Password))
		}
		return handleFunc(c)
	}
}
