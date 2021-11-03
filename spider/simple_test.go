package spider

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
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
		{
			name: "处理iframe",
			args: args{c: core.NewContext(nil, nil, nil)},
			stub: func(c *core.Context) {
				req, _ := http.NewRequest(http.MethodGet, "https://baidu.com/z.html", nil)
				c.SetRequest(req)
				c.SetResponse(core.NewResponse(
					&http.Response{Body: ioutil.NopCloser(
						bytes.NewBufferString(
							`<!DOCTYPE html>
							<html>
							<body>
							<h1>The iframe element</h1>
							<iframe src="https://baidu.com/z.html" title="W3Schools Free Online Web Tutorials">
							</iframe>
							</body>
							</html>`))}))
			},
			wantUrls: []string{"https://baidu.com/z.html"},
		},
		{
			name: "处理iframe和a标签",
			args: args{c: core.NewContext(nil, nil, nil)},
			stub: func(c *core.Context) {
				req, _ := http.NewRequest(http.MethodGet, "https://baidu.com/z.html", nil)
				c.SetRequest(req)
				c.SetResponse(core.NewResponse(
					&http.Response{Body: ioutil.NopCloser(
						bytes.NewBufferString(
							`<!DOCTYPE html>
							<html>
							<body>
							<h1>The iframe element</h1>
							<a href="https://baidu.com/x"></a>
							<iframe src="https://baidu.com/z.html" title="W3Schools Free Online Web Tutorials"></iframe>
							</body>
							</html>`))}))
			},
			wantUrls: []string{"https://baidu.com/x", "https://baidu.com/z.html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub(tt.args.c)
			spider := NewSimpleSpider(nil)
			gotUrls, err := spider.ExtractURL(tt.args.c)
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

func TestURLJoin(t *testing.T) {
	type args struct {
		base string
		href string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "绝对地址",
			args: args{base: "https://baidu.com", href: "https://google.com"},
			want: "https://google.com",
		},
		{
			name: "绝对路径",
			args: args{base: "https://baidu.com/a/b/c", href: "/c"},
			want: "https://baidu.com/c",
		},
		{
			name: "相对路径",
			args: args{base: "https://baidu.com/a/b/c", href: "./d"},
			want: "https://baidu.com/a/b/d",
		},
		{
			name: "相对路径带参数",
			args: args{base: "https://baidu.com/a/b/c", href: "d?a=b"},
			want: "https://baidu.com/a/b/d?a=b",
		},
		{
			name: "相对路径带参数",
			args: args{base: "https://baidu.com/a/b/c?a=x", href: "d?a=b"},
			want: "https://baidu.com/a/b/d?a=b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, _ := url.Parse(tt.args.base)
			got, err := URLJoin(base, tt.args.href)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLJoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("URLJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
