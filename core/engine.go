package core

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"

	"icode.baidu.com/baidu/go-lib/log/log4go"
)

// Echo is the top-level framework instance.
type Engine struct {
	wg          sync.WaitGroup
	middlewares []Middleware
	scheduler   Scheduler
	logger      log4go.Logger
	c           chan struct{}
}

// New creates an instance of Engine.
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
	for _, m := range e.middlewares {
		starter, ok := m.(interface{ Starter })
		if !ok {
			continue
		}
		starter.Start(e)
	}

	logs := make([]string, len(e.middlewares))
	for i, m := range e.middlewares {
		logs[i] = getTypeName(m)
	}

	e.Logger().Info("Use spider middlewares: \n[%s]", strings.Join(logs, ",\n"))

	go func() {
		for {
			task := e.scheduler.BlockPop()
			if task == nil {
				return
			}
			e.c <- struct{}{}
			go func(task *Context) {
				err := task.handler(task)
				// 考虑在这里增加扩展处理
				_ = err
				<-e.c
				e.wg.Done()
			}(task)
		}
	}()
	spider.StartRequests(e)
	c := make(chan struct{})
	go func() {
		e.Shutdown()
		c <- struct{}{}
	}()
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
	for _, m := range e.middlewares {
		closer, ok := m.(Closer)
		if !ok {
			continue
		}
		closer.Close(e)
	}
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
