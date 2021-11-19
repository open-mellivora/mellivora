package middlewares

import (
	"fmt"
	"github.com/open-mellivora/mellivora"
	"sort"
	"strings"

	collector2 "github.com/open-mellivora/mellivora/library/collector"
)

// StatsCollector for collecting scraping stats
type StatsCollector struct {
	groupCollector *collector2.GroupCollector
}

// StatsCollectorConfig defines the config for StatsCollector middleware.
type StatsCollectorConfig struct{}

// DefaultStatsCollectorConfig is the default StatsCollector middleware config.
var DefaultStatsCollectorConfig = StatsCollectorConfig{}

// NewStatsCollectorWithConfig returns a StatsCollector middleware with config.
// See: `NewStatsCollector()`.
func NewStatsCollectorWithConfig(config StatsCollectorConfig) *StatsCollector {
	return &StatsCollector{
		groupCollector: collector2.NewGroupCollect(),
	}
}

// NewStatsCollector returns a StatsCollector instance
func NewStatsCollector() *StatsCollector {
	return NewStatsCollectorWithConfig(DefaultStatsCollectorConfig)
}

// Next implement mellivora.StatsCollector.Next
func (s *StatsCollector) Next(handleFunc mellivora.MiddlewareHandlerFunc)  mellivora.MiddlewareHandlerFunc {
	return func(c *mellivora.Context) error {
		domain := fmt.Sprint(c.GetRequest().URL.Host)
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

// Close implement mellivora.Closer
func (s *StatsCollector) Close(c *mellivora.Engine) {
	sorties := sortKVS{}
	s.groupCollector.Range(func(s string, c *collector2.Collector) bool {
		sorties = append(sorties, kv{key: s, value: c.Count()})
		return true
	})
	sort.Sort(sorties)
	msgs := []string{"Dumping Spider Stats:"}
	for _, item := range sorties {
		msgs = append(msgs, fmt.Sprintf("'%v': %v", item.key, item.value))
	}
	c.Logger().Sugar().Infof(strings.Join(msgs, "\n"))
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
