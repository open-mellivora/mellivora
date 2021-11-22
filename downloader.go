package mellivora

import (
	"net/http"
)

type Downloader struct{}

func (d *Downloader) Next(handler MiddlewareFunc) MiddlewareFunc {
	return func(c *Context) (err error) {
		if c.Response != nil {
			return nil
		}

		var response *http.Response
		if c.roundTripper == nil {
			c.roundTripper = http.DefaultTransport
		}

		if response, err = d.download(&http.Client{
			Transport: c.roundTripper,
		},
			c.GetRequest()); err != nil {
			return
		}

		c.SetResponse(NewResponse(response))
		return handler(c)
	}
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) download(client *http.Client, req *http.Request) (
	response *http.Response, err error) {

	if client == nil {
		client = &http.Client{}
	}
	return client.Do(req)
}
