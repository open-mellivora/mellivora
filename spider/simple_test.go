package spider

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

func TestExtractURL(t *testing.T) {
	type args struct {
		c *core.Context
	}
	tests := []struct {
		name     string
		args     args
		wantUrls []string
		wantErr  bool
		stub     func(c *core.Context)
	}{
		{
			name: "不同层级a标签",
			args: args{c: core.NewContext(nil, nil, nil)},
			stub: func(c *core.Context) {
				req, _ := http.NewRequest(http.MethodGet, "https://baidu.com/z.html", nil)
				c.SetRequest(req)
				c.SetResponse(core.NewResponse(&http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(
					`<html><body><a href="/x"></a><div><a href="https://baidu.com/y.html"></a>
					</div></body></html>`))}))
			},
			wantUrls: []string{"https://baidu.com/x", "https://baidu.com/y.html"},
		},
		{
			name: "过滤非当前域名链接",
			args: args{c: core.NewContext(nil, nil, nil)},
			stub: func(c *core.Context) {
				req, _ := http.NewRequest(http.MethodGet, "https://baidu.com/z.html", nil)
				c.SetRequest(req)
				c.SetResponse(core.NewResponse(
					&http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(
						`<html><body><a href="/x"></a><div><a style='' href="https://zhihu.com/y.html"></a>
						</div></body></html>`))}))
			},
			wantUrls: []string{"https://baidu.com/x"},
		},
		{
			name: "过滤异常链接",
			args: args{c: core.NewContext(nil, nil, nil)},
			stub: func(c *core.Context) {
				req, _ := http.NewRequest(http.MethodGet, "https://baidu.com/z.html", nil)
				c.SetRequest(req)
				c.SetResponse(core.NewResponse(
					&http.Response{Body: ioutil.NopCloser(
						bytes.NewBufferString(
							`<html><body><a href="javascript:;"></a><div>
							<a href="https://zhihu.com:65536/y.html"></a></div></body></html>`))}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub(tt.args.c)
			gotUrls, err := ExtractURL(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUrls, tt.wantUrls) {
				t.Errorf("ExtractURL() gotUrls = %v, want %v", gotUrls, tt.wantUrls)
			}
		})
	}
}
