package downlimiter

import (
	"time"

	"golang.org/x/net/context"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
)

type Middleware struct {
	cfg                         *Config
	concurrencyLimiter          *ConcurrencyLimiter
	concurrencyPerDomainLimiter *ConcurrencyGroupLimiter
	downloadDelayPerDomain      *DelayGroupLimiter
}

type Config struct {
	ConcurrentRequests          int64
	ConcurrentRequestsPerDomain int64
	DownloadDelayPerDomain      time.Duration
	Timeout                     time.Duration
	MaxDepth                    int64
}

var DefaultConfig = &Config{
	ConcurrentRequests:          10,
	ConcurrentRequestsPerDomain: 5,
	DownloadDelayPerDomain:      time.Second,
	Timeout:                     1,
}

func NewMiddleware(config *Config) *Middleware {
	if config == nil {
		config = DefaultConfig
	}
	if config.ConcurrentRequestsPerDomain == 0 {
		config.ConcurrentRequestsPerDomain = 1
	}
	if config.ConcurrentRequests == 0 {
		config.ConcurrentRequests = 1
	}
	m := &Middleware{
		cfg:                         config,
		concurrencyLimiter:          NewConcurrencyLimiter(config.ConcurrentRequests),
		concurrencyPerDomainLimiter: NewConcurrencyGroupLimiter(config.ConcurrentRequestsPerDomain),
		downloadDelayPerDomain:      NewDelayGroupLimiter(config.DownloadDelayPerDomain),
	}
	return m
}

func (m *Middleware) Next(handleFunc core.HandleFunc) core.HandleFunc {
	return func(c *core.Context) (err error) {
		if c.GetDepth() > m.cfg.MaxDepth {
			return nil
		}
		req := c.GetRequest()
		domain := req.URL.Host
		m.concurrencyPerDomainLimiter.Wait(domain)
		m.concurrencyLimiter.Wait()
		m.downloadDelayPerDomain.Wait(domain)
		defer func() {
			m.concurrencyLimiter.Done()
			m.concurrencyPerDomainLimiter.Done(domain)
			m.downloadDelayPerDomain.Reset(domain)
		}()
		if m.cfg.Timeout != 0 {
			ctx := req.Context()
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, m.cfg.Timeout)
			defer cancel()
			c.SetRequest(req.WithContext(ctx))
		}
		err = handleFunc(c)
		return err
	}
}
