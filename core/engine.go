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

type Engine struct {
	wg          sync.WaitGroup
	middlewares []Middleware
	scheduler   Scheduler
	logger      log4go.Logger
	c           chan struct{}
}

func NewEngine() *Engine {
	core := &Engine{
		wg:        sync.WaitGroup{},
		scheduler: NewLifoScheduler(),
		logger:    log4go.NewDefaultLogger(log4go.INFO),
		c:         make(chan struct{}, 128),
	}
	return core
}

// Use 添加中间件，优先级从左往右
func (e *Engine) Use(middlewares ...Middleware) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e *Engine) SetLogger(l log4go.Logger) {
	e.logger = l
}

func (e *Engine) Logger() log4go.Logger {
	return e.logger
}

// applyMiddleware 中间件组合
// middlewares优先级从左往右
func (e *Engine) applyMiddleware(h HandleFunc, middlewares ...Middleware) HandleFunc {
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

// Run 阻塞运行spider
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
		logs[i] = GetTypeName(m)
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

// Shutdown 捕获退出信号，安全退出
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

func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}
