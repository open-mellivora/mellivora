package saver

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/valyala/bytebufferpool"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

type Middleware struct {
	r            *regexp.Regexp
	config       *Config
	OpenFileFunc func(filename string) (io.WriteCloser, error)
	bytesPool    *bytebufferpool.Pool
}

// dirCreate checks and creates dir if nonexist
func dirCreate(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(dir, 0o777)
		if err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	Dir          string
	Pattern      string
	OpenFileFunc func(filename string) (io.WriteCloser, error)
}

func NewMiddleware(config *Config) (mid *Middleware, err error) {
	var r *regexp.Regexp
	if r, err = regexp.Compile(config.Pattern); err != nil {
		return
	}
	if err = dirCreate(config.Dir); err != nil {
		return
	}

	m := &Middleware{
		config:       config,
		r:            r,
		OpenFileFunc: config.OpenFileFunc,
		bytesPool:    &bytebufferpool.Pool{},
	}
	if m.OpenFileFunc == nil {
		m.OpenFileFunc = func(filename string) (io.WriteCloser, error) {
			return os.Create(filename)
		}
	}
	return m, err
}

func (m *Middleware) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) (err error) {
		u := c.GetRequest().URL
		if err = handleFunc(c); err != nil {
			return
		}
		if c.Response == nil {
			return handleFunc(c)
		}

		if !m.r.MatchString(u.String()) {
			return handleFunc(c)
		}

		filename := url.QueryEscape(u.String())
		path := filepath.Join(m.config.Dir, filename)
		var f io.WriteCloser
		f, err = m.OpenFileFunc(path)
		if err != nil {
			return err
		}
		defer f.Close()

		buf := m.bytesPool.Get()
		defer m.bytesPool.Put(buf)
		newReader := io.TeeReader(c.Response.Body, buf)
		c.Response.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		if _, err = io.Copy(f, newReader); err != nil {
			return
		}
		return
	}
}
