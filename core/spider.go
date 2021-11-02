package core

type Spider interface {
	// StartRequests spider启动函数
	StartRequests(c *Context) error
}
