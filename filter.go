package mellivora

type Filter interface {
	Exist(c *Context) bool
	Add(c *Context)
}
