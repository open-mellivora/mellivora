package mellivora

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/open-mellivora/mellivora/library/limter"
)

// Engine is the top-level framework instance.
type Engine struct {
	contextSerializer  *ContextSerializer
	wg                 sync.WaitGroup
	middlewares        []Middleware
	scheduler          Scheduler
	logger             *zap.Logger
	concurrencyLimiter *limter.ConcurrencyLimiter
	filter             Filter
}

// NewEngine creates an instance of Engine.
func NewEngine(concurrency uint64) *Engine {
	if concurrency <= 0 {
		concurrency = 1
	}
	core := &Engine{
		wg:                 sync.WaitGroup{},
		scheduler:          NewLifoScheduler(),
		contextSerializer:  NewContextSerializer(),
		concurrencyLimiter: limter.NewConcurrencyLimiter(concurrency),
		filter:             NewBloomFilter(),
		logger:             zap.NewExample(),
	}
	return core
}

// Use adds middleware to the chain which is run before spider.
func (e *Engine) Use(middlewares ...Middleware) {
	e.middlewares = append(e.middlewares, middlewares...)
}

// SetLogger sets `*log.Logger`.
func (e *Engine) SetLogger(l *zap.Logger) {
	e.logger = l
}

// Logger returns `*log.Logger`.
func (e *Engine) Logger() *zap.Logger {
	return e.logger
}

func (e *Engine) applyMiddleware(middlewares ...Middleware) MiddlewareFunc {
	h := func(c *Context) error {
		return nil
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		handle := m.Next(h)
		h = func(c *Context) error {
			err := handle(c)
			return err
		}
	}
	return h
}

func (e *Engine) runC(c *Context, middlewareFunc MiddlewareFunc) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger().Error("Recover", zap.Any("recover", r))
		}
	}()

	if err := middlewareFunc(c); err != nil {
		return
	}

	if c.Response == nil {
		return
	}
	var task Task
	if task = c.handler(c); task == nil {
		e.Logger().Warn("Empty Task", zap.String("url", c.GetRequest().URL.String()))
	}
	for i := 0; i < len(task.Requests()); i++ {
		e.request(c, task.Requests()[i], task.Handler(), task.RequestOptions()...)
	}
}

// Run run a spider.
func (e *Engine) Run(spider Spider) {
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		m := e.middlewares[i]
		starter, ok := m.(interface{ Runnable })
		if !ok {
			continue
		}
		e.Logger().Info("Middleware Start", zap.String("Name", getTypeName(m)))
		starter.Start()
	}

	names := make([]string, len(e.middlewares))
	for i, m := range e.middlewares {
		names[i] = getTypeName(m)
	}

	e.Logger().Info("Use spider middlewares", zap.Strings("names", names))

	middlewareFunc := e.applyMiddleware(append(e.middlewares, NewDownloader())...)

	go func() {
		for {
			contextText := e.scheduler.Pop()
			if contextText == nil {
				continue
			}
			var c *Context
			var err error
			if c, err = e.contextSerializer.Unmarshal(contextText); err != nil {
				e.Logger().Error("unmarshal c error", zap.Error(err))
				e.wg.Done()
				continue
			}
			c.core = e

			e.concurrencyLimiter.Wait()
			go func(c *Context) {
				e.filter.Add(c)
				e.runC(c, middlewareFunc)
				e.concurrencyLimiter.Done()
				e.wg.Done()
			}(c)
		}
	}()
	var task Task
	if task = spider.StartRequests(); task == nil {
		e.Logger().Warn("empty task")
	}

	for _, r := range task.Requests() {
		e.request(nil, r, task.Handler(), task.RequestOptions()...)
	}

	c := make(chan struct{})
	// wait shutdown
	go func() {
		e.Shutdown()
		c <- struct{}{}
	}()
	// wait c done
	go func() {
		e.wg.Wait()
		c <- struct{}{}
	}()
	<-c
	e.Close()
}

// Close immediately stops the server.
func (e *Engine) Close() {
	e.scheduler.Close()
	e.Logger().Info("Scheduler Closed")
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		m := e.middlewares[i]
		closer, ok := m.(Closable)
		if !ok {
			continue
		}
		closer.Close()
		e.Logger().Info("Middleware Closed", zap.String("name", getTypeName(m)))
	}
	e.Logger().Info("Engine Closed")
	e.Logger().Sync()
}

// Shutdown stops the server gracefully.
func (e *Engine) Shutdown() {
	quit := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(quit, os.Interrupt)
	go func() {
		sig := <-quit
		e.Logger().Info("hutting down, caused by", zap.String("signal", sig.String()))
		close(done)
	}()
	<-done
	e.Logger().Info("Graceful shutdown")
}

func getTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

// request create a request
func (e *Engine) request(preCtx *Context, r *http.Request, handler HandleFunc,
	options ...RequestOptionsFunc) (err error) {
	if r.URL.Scheme == "" {
		err = errors.New("empty scheme")
		return
	}

	ctx := NewContext(e, r, handler)

	opt := NewRequestOptions()
	for _, optFunc := range options {
		optFunc(opt)
	}
	ctx.setter = opt.setter
	if preCtx != nil {
		ctx.SetDepth(preCtx.GetDepth() + 1)
	}

	if e.filter.Exist(ctx) {
		return
	}

	e.wg.Add(1)
	var bs []byte
	if bs, err = e.contextSerializer.Marshal(ctx); err != nil {
		return
	}
	e.scheduler.Push(bs)
	return
}
