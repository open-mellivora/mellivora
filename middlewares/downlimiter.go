package middlewares

import (
	"time"

	"golang.org/x/net/context"

	"github.com/open-mellivora/mellivora"
	"github.com/open-mellivora/mellivora/library/limter"
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
	ConcurrentRequestsPerDomain uint64
	// DownloadDelayPerDomain 每个域名下请求的延时
	DownloadDelayPerDomain time.Duration
	// Timeout 请求超时
	Timeout time.Duration
}

// DefaultDownLimiterConfig is the default DownLimiter middleware config.
var DefaultDownLimiterConfig = DownLimiterConfig{
	ConcurrentRequestsPerDomain: 1024,
	DownloadDelayPerDomain:      0,
	Timeout:                     3 * time.Second,
}

// NewDownLimiterWithConfig returns a DownLimiter middleware with config.
// See: `DownLimiter()`.
func NewDownLimiterWithConfig(config DownLimiterConfig) *DownLimiter {
	if config.ConcurrentRequestsPerDomain == 0 {
		config.ConcurrentRequestsPerDomain =
			DefaultDownLimiterConfig.ConcurrentRequestsPerDomain
	}
	m := &DownLimiter{
		config: config,
		concurrencyPerDomainLimiter: limter.NewConcurrencyGroupLimiter(
			config.ConcurrentRequestsPerDomain),
		downloadDelayPerDomain: limter.NewDelayGroupLimiter(
			config.DownloadDelayPerDomain),
	}
	return m
}

// NewDownLimiter returns a DownLimiter instance
func NewDownLimiter() *DownLimiter {
	return NewDownLimiterWithConfig(DefaultDownLimiterConfig)
}

// Next implement mellivora.Middleware.Next
func (m *DownLimiter) Next(handleFunc mellivora.MiddlewareFunc) mellivora.MiddlewareFunc {
	return func(c *mellivora.Context) (err error) {
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
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, m.config.Timeout)
			defer cancel()
			c.SetRequest(req.WithContext(ctx))
		}
		err = handleFunc(c)
		return err
	}
}
