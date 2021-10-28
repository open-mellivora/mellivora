package middlewares

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"golang.org/x/text/encoding/simplifiedchinese"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

func TestDecoder(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	c := core.NewContext(nil, req, nil)
	decoder := NewDecoder()
	t.Run("解码gbk成功", func(t *testing.T) {
		bs, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(`<html>
<!DOCTYPE html>
<html>
<head>
    <meta charset="gbk">
    <title>Document</title>
</head>
<body>
    <span>你好</span>
</body>
</html>
`))
		resp := &http.Response{
			StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer(bs)),
			Header: map[string][]string{
				"Content-Type": {"Content-type:text/html;charset=gbk"},
			},
		}
		c.SetResponse(core.NewResponse(resp))
		err := decoder.Next(func(c *core.Context) error {
			return nil
		})(c)
		assert.Equal(t, err, nil)
		str, err := c.String()
		assert.Equal(t, err, nil)
		assert.Equal(t, strings.Contains(str, "你好"), true)
	})
	t.Run("解码utf-8编码成功", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewBuffer([]byte(`<html>
<!DOCTYPE html>
<html>
<head>
    <title>Document</title>
</head>
<body>
    <span>你好</span>
</body>
</html>
`))), Header: map[string][]string{
				"Content-Type": {"Content-type:text/html"},
			},
		}
		c.SetResponse(core.NewResponse(resp))
		err := decoder.Next(func(c *core.Context) error {
			return nil
		})(c)
		assert.Equal(t, err, nil)
		str, err := c.String()
		assert.Equal(t, err, nil)
		assert.Equal(t, strings.Contains(str, "你好"), true)
	})

	t.Run("外层失败", func(t *testing.T) {
		err := decoder.Next(func(c *core.Context) error {
			return net.ErrClosed
		})(c)
		assert.NotEqual(t, err, nil)
	})
}
