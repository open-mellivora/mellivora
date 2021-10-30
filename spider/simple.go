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

func (s *SimpleSpider) StartRequests(c *core.Context) error {
	for i := 0; i < len(s.urls); i++ {
		if err := c.Get(s.urls[i], s.Parse); err != nil {
			continue
		}
	}
	return nil
}

func NewURL(base *url.URL, href string) (*url.URL, error) {
	hrefURL, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(hrefURL), nil
}

func ExtractURL(c *core.Context) (urls []string, err error) {
	tokenizer := c.Tokenizer()

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

	return c.Gets(urls, s.Parse)
}
