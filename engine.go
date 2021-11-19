package mellivora

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"

	 "github.com/open-mellivora/mellivora/library/limter"
	"go.uber.org/zap"
)

// Engine is the top-level framework instance.
type Engine struct {
	contextSerializer  *ContextSerializer
	wg                 sync.WaitGroup
	middlewares        []Middleware
	scheduler          Scheduler
	logger             *zap.Logger
	concurrencyLimiter *limter.ConcurrencyLimiter
}

// NewEngine creates an instance of Engine.
func NewEngine(concurrency int64) *Engine {
	if concurrency <= 0 {
		concurrency = 1
	}
	core := &Engine{
		wg:                 sync.WaitGroup{},
		scheduler:          NewLifoScheduler(),
		logger:             zap.NewExample(),
		contextSerializer:  NewContextSerializer(),
		concurrencyLimiter: limter.NewConcurrencyLimiter(concurrency),
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

func (e *Engine) applyMiddleware(middlewares ...Middleware) MiddlewareHandlerFunc {
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

func (e *Engine) runc(c *Context, middlewareFunc MiddlewareHandlerFunc) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger().Sugar().Errorf("Recovery: %+v", r)
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
		e.Logger().Sugar().Warnf("Empty %s Task", c.GetRequest().URL.String())
	}
	for i:=0;i<len(task.Requests());i++{
		e.request(c,task.Requests()[i],task.Handler(),task.RequestOptions()...)
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
		e.Logger().Sugar().Infof("Middleware %s Start", getTypeName(m))
		starter.Start()
	}

	logs := make([]string, len(e.middlewares))
	for i, m := range e.middlewares {
		logs[i] = getTypeName(m)
	}

	e.Logger().Sugar().Infof("Use spider middlewares: \n[%s]", strings.Join(logs, ",\n"))

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
				e.Logger().Sugar().Errorf("unmarshal c error,err:%s", err.Error())
				continue
			}
			c.core = e

			e.concurrencyLimiter.Wait()
			go func(c *Context) {
				e.runc(c, middlewareFunc)
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
		e.Logger().Sugar().Infof("Middleware %s Closed", getTypeName(m))
	}
	e.Logger().Info("Engine Closed")
}

// Shutdown stops the server gracefully.
func (e *Engine) Shutdown() {
	quit := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(quit, os.Interrupt)
	go func() {
		sig := <-quit
		log.Println("Shutting down, caused by", sig)
		close(done)
	}()
	<-done
	log.Println("Graceful shutdown.")
}

func getTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

// request create a request
func (e *Engine) request(preCtx *Context, r *http.Request, handler HandleFunc,
	options ...RequestOptionsFunc) (err error) {

	ctx := NewContext(e, r, handler)

	opt := NewRequestOptions()
	for _, optFunc := range options {
		optFunc(opt)
	}
	ctx.setter = opt.setter
	if preCtx != nil {
		ctx.SetHTTPClient(preCtx.httpClient)
	}
	e.wg.Add(1)
	var bs []byte
	if bs, err = e.contextSerializer.Marshal(ctx); err != nil {
		return
	}
	e.scheduler.Push(bs)
	return
}
