package main

import (
	"github.com/open-mellivora/mellivora"
	"github.com/open-mellivora/mellivora/middlewares"
	"golang.org/x/net/html"
	"net/url"
)

func NewSimpleSpider(urls []string) *SimpleSpider {
	return &SimpleSpider{
		urls: urls,
		matchs: map[string]string{
			"a":      "href",
			"iframe": "src",
		},
	}
}

type SimpleSpider struct {
	urls   []string
	matchs map[string]string // map[tag]attr
}

// StartRequests implement core.Spider.StartRequests
func (s *SimpleSpider) StartRequests() mellivora.Task {
	task,_:= mellivora.Gets(s.urls,s.Parse)
	return task
}

// URLJoin Construct a full absolute URL by combining a base URL with another URL
func URLJoin(base *url.URL, href string) (*url.URL, error) {
	hrefURL, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(hrefURL), nil
}

// ExtractURL Extract links from iframe and a
func (s *SimpleSpider) ExtractURL(c *mellivora.Context) (urls []string, err error) {
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
		attr, has := s.matchs[string(tn)]
		if !has {
			continue
		}

		more := true
		var key, value []byte
		for more {
			key, value, more = tokenizer.TagAttr()
			if string(key) != attr {
				continue
			}

			v := string(value)
			if v == "javascript:;" || v == "javascript:void(0)" {
				break
			}

			newURL, err := URLJoin(c.GetRequest().URL, v)
			if err != nil {
				break
			}

			// Only allow the current domain
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

func (s *SimpleSpider) Parse(c *mellivora.Context) (task mellivora.Task) {
	var urls []string
	var err error
	if urls, err = s.ExtractURL(c); err != nil {
		return
	}

	task,err= mellivora.Gets(urls,s.Parse)
	return task
}

func main(){
	spider:=NewSimpleSpider([]string{"https://www.sina.com.cn/"})
	e:=mellivora.NewEngine(1)
	e.Use(middlewares.NewLogging())
	e.Run(spider)
}