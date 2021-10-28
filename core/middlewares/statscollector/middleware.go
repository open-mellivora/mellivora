package statscollector

import (
	"fmt"
	"sort"
	"strings"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

// Middleware 状态收集中间件
type Middleware struct {
	groupCollector *GroupCollector
}

func NewMiddleware() *Middleware {
	return &Middleware{
		groupCollector: NewGroupCollect(),
	}
}

// Next implement core.Middleware.Next
func (s *Middleware) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) error {
		domain := fmt.Sprint(c.GetRequest().URL.Host)
		s.groupCollector.Add("request_bytes", c.GetRequest().ContentLength)
		s.groupCollector.Add("request_count", 1)
		err := handleFunc(c)
		s.groupCollector.Add("response_count", 1)
		if c.Response != nil {
			s.groupCollector.Add(fmt.Sprint("response_count/", domain), 1)
			s.groupCollector.Add(fmt.Sprint("response_count/", domain, "/", c.Response.StatusCode), 1)
			s.groupCollector.Add(fmt.Sprint("response_count/", c.Response.StatusCode), 1)
		}
		if err != nil {
			s.groupCollector.Add("error_count", 1)
			s.groupCollector.Add(fmt.Sprint("error_count/", domain), 1)
		}
		return err
	}
}

// Close implement core.Closer
func (s *Middleware) Close(c *core.Engine) {
	sortkvs := sortKVS{}
	s.groupCollector.Range(func(key, value interface{}) bool {
		sortkvs = append(sortkvs, kv{key: fmt.Sprint(key), value: value.(*Collector).i})
		return true
	})
	sort.Sort(sortkvs)
	msgs := []string{"Dumping Spider Stats:"}
	for _, item := range sortkvs {
		msgs = append(msgs, fmt.Sprintf("'%v': %v", item.key, item.value))
	}
	c.Logger().Info(strings.Join(msgs, "\n"))
}

type kv struct {
	key   string
	value interface{}
}

type sortKVS []kv

func (s sortKVS) Len() int {
	return len(s)
}

func (s sortKVS) Less(i, j int) bool {
	leftKey := s[i].key
	rightKey := s[j].key
	leftLen := len(leftKey)
	rightLen := len(rightKey)
	for i := 0; i < leftLen && i < rightLen; i++ {
		if leftKey[i] < rightKey[i] {
			return true
		}
		if leftKey[i] > rightKey[i] {
			return false
		}
		continue
	}
	if leftLen < rightLen {
		return true
	}

	return leftKey < rightKey
}

func (s sortKVS) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
