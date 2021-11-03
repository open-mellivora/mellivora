package core

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"

	"icode.baidu.com/baidu/go-lib/log/log4go"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/library/limter"
)

// Engine is the top-level framework instance.
type Engine struct {
	contextSerializer  *ContextSerializer
	wg                 sync.WaitGroup
	middlewares        []Middleware
	scheduler          Scheduler
	logger             log4go.Logger
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
		logger:             log4go.NewDefaultLogger(log4go.INFO),
		contextSerializer:  NewContextSerializer(),
		concurrencyLimiter: limter.NewConcurrencyLimiter(concurrency),
	}
	return core
}

// Use adds middleware to the chain which is run before spider.
func (e *Engine) Use(middlewares ...Middleware) {
	e.middlewares = append(e.middlewares, middlewares...)
}

// SetLogger sets `log4go.Logger`.
func (e *Engine) SetLogger(l log4go.Logger) {
	e.logger = l
}

// Logger returns `log4go.Logger`.
func (e *Engine) Logger() log4go.Logger {
	return e.logger
}

func (e *Engine) applyMiddleware(middlewares ...Middleware) HandlerFunc {
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

func (e *Engine) runTask(task *Context, middlewareFunc HandlerFunc) {
	if err := middlewareFunc(task); err != nil {
		return
	}

	if task.Response == nil {
		return
	}

	if err := task.handler(task); err != nil {
		e.Logger().Error("Parse %s error", task.GetRequest().URL.String())
	}
}

// Run run a spider.
func (e *Engine) Run(spider Spider) {
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		m := e.middlewares[i]
		starter, ok := m.(interface{ Starter })
		if !ok {
			continue
		}
		e.Logger().Info("Middleware %s Start", getTypeName(m))
		starter.Start(e)
	}

	logs := make([]string, len(e.middlewares))
	for i, m := range e.middlewares {
		logs[i] = getTypeName(m)
	}

	e.Logger().Info("Use spider middlewares: \n[%s]", strings.Join(logs, ",\n"))

	middlewareFunc := e.applyMiddleware(append(e.middlewares, NewDownloader())...)

	go func() {
		for {
			taskText := e.scheduler.Pop()
			if taskText == nil {
				continue
			}
			var task *Context
			var err error
			if task, err = e.contextSerializer.Unmarshal(taskText); err != nil {
				e.Logger().Error("unmarshal task error,err:%s", err.Error())
				continue
			}
			task.core = e

			e.concurrencyLimiter.Wait()
			go func(task *Context) {
				e.runTask(task, middlewareFunc)
				e.concurrencyLimiter.Done()
				e.wg.Done()
			}(task)
		}
	}()

	ctx := NewContext(e, nil, nil)
	if err := spider.StartRequests(ctx); err != nil {
		e.Logger().Error("start requests error err:%s", err.Error())
	}

	c := make(chan struct{})
	// wait shutdown
	go func() {
		e.Shutdown()
		c <- struct{}{}
	}()
	// wait task done
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
		closer, ok := m.(Closer)
		if !ok {
			continue
		}
		closer.Close(e)
		e.Logger().Info("Middleware %s Closed", getTypeName(m))
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
func (e *Engine) request(preCtx *Context, r *http.Request, handler HandlerFunc,
	options ...RequestOptionsFunc) (err error) {

	ctx := NewContext(e, r, handler)

	opt := NewRequestOptions()
	for _, optFunc := range options {
		optFunc(opt)
	}
	ctx.setter = opt.setter
	ctx.SetHTTPClient(preCtx.httpClient)
	e.wg.Add(1)
	var bs []byte
	if bs, err = e.contextSerializer.Marshal(ctx); err != nil {
		return
	}
	e.scheduler.Push(bs)
	return
}
