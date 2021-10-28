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
)

// Engine is the top-level framework instance.
type Engine struct {
	wg          sync.WaitGroup
	middlewares []Middleware
	scheduler   Scheduler
	logger      log4go.Logger
	c           chan struct{}
}

// NewEngine creates an instance of Engine.
func NewEngine() *Engine {
	core := &Engine{
		wg:        sync.WaitGroup{},
		scheduler: NewLifoScheduler(),
		logger:    log4go.NewDefaultLogger(log4go.INFO),
		c:         make(chan struct{}, 128),
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

func (e *Engine) applyMiddleware(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
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
	go func() {
		for {
			task := e.scheduler.Pop()
			if task == nil {
				continue
			}
			e.c <- struct{}{}
			go func(task *Context) {
				err := task.handler(task)
				// 考虑在这里增加扩展处理error
				_ = err
				<-e.c
				e.wg.Done()
			}(task)
		}
	}()
	spider.StartRequests(e)
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
		log.Println("Shutting down, caused by ", sig)
		close(done)
	}()
	<-done
	log.Println("Graceful shutdown.")
}

func getTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

// Get create a GET request
func (e *Engine) Get(url string, handler HandlerFunc, options ...RequestOptionsFunc) (err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	e.Request(req, handler, options...)
	return
}

// Request create a request
func (e *Engine) Request(r *http.Request, handler HandlerFunc, options ...RequestOptionsFunc) {
	middlewares := append(e.middlewares, NewDownloader())

	middlewareHandler := e.applyMiddleware(func(c *Context) error {
		return nil
	}, middlewares...)

	ctx := NewContext(e, r, func(c *Context) error {
		if err := middlewareHandler(c); err != nil {
			return err
		}
		// 过滤等正常情况导致Response无数据
		if c.Response == nil {
			return nil
		}
		return handler(c)
	})

	opt := RequestOptions{}
	for _, optFunc := range options {
		optFunc(&opt)
	}

	if opt.PreContext != nil {
		ctx.SetDepth(opt.PreContext.GetDepth() + 1)
	}

	e.wg.Add(1)
	e.scheduler.Push(ctx)
}
