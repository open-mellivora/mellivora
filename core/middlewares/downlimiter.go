package middlewares

import (
	"time"

	"golang.org/x/net/context"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/library/limter"
)

// DownLimiter for limit downloader
type DownLimiter struct {
	config                      DownLimiterConfig
	concurrencyPerDomainLimiter *limter.ConcurrencyGroupLimiter
	downloadDelayPerDomain      *limter.DelayGroupLimiter
}

// DownLimiterConfig  defines the config for DownLimiter middleware.
type DownLimiterConfig struct {
	// ConcurrentRequestsPerDomain 每个域名下请求并行限制
	ConcurrentRequestsPerDomain int64
	// DownloadDelayPerDomain 每个域名下请求的延时
	DownloadDelayPerDomain time.Duration
	// Timeout 请求超时
	Timeout  time.Duration
	MaxDepth int64
}

// DefaultDownLimiterConfig is the default DownLimiter middleware config.
var DefaultDownLimiterConfig = DownLimiterConfig{
	ConcurrentRequestsPerDomain: 5,
	DownloadDelayPerDomain:      time.Second,
	Timeout:                     1,
}

// NewDownLimiterWithConfig returns a DownLimiter middleware with config.
// See: `DownLimiter()`.
func NewDownLimiterWithConfig(config DownLimiterConfig) *DownLimiter {
	if config.ConcurrentRequestsPerDomain == 0 {
		config.ConcurrentRequestsPerDomain = 1
	}

	m := &DownLimiter{
		config: config,
		concurrencyPerDomainLimiter: limter.NewConcurrencyGroupLimiter(
			config.ConcurrentRequestsPerDomain),
		downloadDelayPerDomain: limter.NewDelayGroupLimiter(config.DownloadDelayPerDomain),
	}
	return m
}

// NewDownLimiter returns a DownLimiter instance
func NewDownLimiter() *DownLimiter {
	return NewDownLimiterWithConfig(DefaultDownLimiterConfig)
}

// Next implement core.Middleware.Next
func (m *DownLimiter) Next(handleFunc core.HandlerFunc) core.HandlerFunc {
	return func(c *core.Context) (err error) {
		if c.GetDepth() > m.config.MaxDepth {
			return nil
		}
		req := c.GetRequest()
		domain := req.URL.Host
		m.concurrencyPerDomainLimiter.Wait(domain)
		m.downloadDelayPerDomain.Wait(domain)
		defer func() {
			m.concurrencyPerDomainLimiter.Done(domain)
			m.downloadDelayPerDomain.Reset(domain)
		}()
		if m.config.Timeout != 0 {
			ctx := req.Context()
			ctx, _ = context.WithTimeout(ctx, m.config.Timeout)
			c.SetRequest(req.WithContext(ctx))
		}
		err = handleFunc(c)
		return err
	}
}
