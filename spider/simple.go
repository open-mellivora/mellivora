package spider

import (
	"net/url"

	"golang.org/x/net/html"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

func NewSimpleSpider(urls []string) *SimpleSpider {
	return &SimpleSpider{
		urls: urls,
	}
}

type SimpleSpider struct {
	urls []string
}

func (s *SimpleSpider) StartRequests(c *core.Engine) {
	for i := 0; i < len(s.urls); i++ {
		if err := c.Get(s.urls[i], s.Parse); err != nil {
			continue
		}
	}
}

func NewURL(base *url.URL, href string) (*url.URL, error) {
	hrefURL, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(hrefURL), nil
}

func ExtractURL(c *core.Context) (urls []string, err error) {
	var tokenizer *html.Tokenizer
	tokenizer, err = c.Tokenizer()
	if err != nil {
		c.Engine().Logger().Warn("tokenizer error,err:%v,url:%s", err, c.GetRequest().URL.String())
		return
	}

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt != html.StartTagToken {
			continue
		}

		tn, _ := tokenizer.TagName()
		if len(tn) != 1 || tn[0] != 'a' {
			continue
		}
		more := true
		var key, value []byte
		for more {
			key, value, more = tokenizer.TagAttr()
			if string(key) != "href" {
				continue
			}

			v := string(value)
			if v == "javascript:;" || v == "javascript:void(0)" {
				break
			}

			newURL, err := NewURL(c.GetRequest().URL, v)
			if err != nil {
				break
			}

			// 限制在当前域名
			if newURL.Host != c.GetRequest().URL.Host {
				break
			}

			urls = append(urls, newURL.String())
			if !more {
				break
			}
		}
	}
	return
}

func (s *SimpleSpider) Parse(c *core.Context) (err error) {
	var urls []string
	if urls, err = ExtractURL(c); err != nil {
		return
	}
	for i := 0; i < len(urls); i++ {
		if err = c.Engine().Get(urls[i], s.Parse, core.WithPreContext(c)); err != nil {
			c.Engine().Logger().Warn("get error,url:%s", urls[i])
		}
	}
	return
}
