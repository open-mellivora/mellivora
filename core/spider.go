package core

type Spider interface {
	// StartRequests generate first requests
	StartRequests(c *Context) error
}
