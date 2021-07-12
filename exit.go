package regia

type Exit interface{ Exit(context *Context) }

type exit struct{}

// Exit Do nothing
func (e exit) Exit(*Context) {}


