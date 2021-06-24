package regia

import (
	"net/http"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

const (
	FilePathParam = "FilePathParam"
	wildFilepath  = "*" + FilePathParam
)

// Engine is a collection of core components of the whole service
type Engine struct {
	// The branch are used to store the handler.
	// All handlers are going to register in the router
	*BluePrint

	// Router is a module used to register handle and distribute request
	Router Router

	// NotFoundHandle replies to the request with an HTTP 404 not found error.
	NotFoundHandle func(context *Context)

	// All requests will be intercepted by Interceptors
	// whatever route matched or not
	Interceptors HandleFuncGroup

	// Starter will run when the service starts
	// and it only run once
	Starters []Starter

	// Warehouse is used to store information
	Warehouse Warehouse

	MultipartMemory int64

	pool sync.Pool

	Abort Exit
}

func (e *Engine) dispatchContext() *Context {
	return &Context{Engine: e}
}

// register all handles to router
func (e *Engine) registerHandle() {
	for method, nodes := range e.methodsTree {
		for _, node := range nodes {
			e.Router.Insert(method, node.path, node.group)
		}
	}
}

// SetNotFoundHandle Setter for Engine.NotFoundHandle
func (e *Engine) SetNotFoundHandle(handle HandleFunc) {
	e.NotFoundHandle = handle
}

// Static Serve static files
func (e *Engine) Static(url, dir string, group ...HandleFunc) {
	if strings.Contains(url, "*") {
		panic("`url` should not have wildcards")
	}
	server := http.FileServer(http.Dir(dir))
	handle := func(context *Context) {
		context.Request.URL.Path = context.Params.Get(FilePathParam).Text()
		server.ServeHTTP(context.ResponseWriter, context.Request)
	}
	group = append(group, handle)
	if !strings.HasSuffix(url, FilePathParam) {
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		url += wildFilepath
	}
	e.Handle(http.MethodGet, url, group...)
}

// AddInterceptors Add interceptor to Engine
func (e *Engine) AddInterceptors(interceptors ...HandleFunc) {
	e.Interceptors = append(e.Interceptors, interceptors...)
}

// AddStarter Add starter to Engine
func (e *Engine) AddStarter(starters ...Starter) {
	e.Starters = append(e.Starters, starters...)
}

// Call all starters of this engine
func (e *Engine) runStarter() {
	for _, starter := range e.Starters {
		starter.Start(e)
	}
}

// Init engine
func (e *Engine) init() {
	e.registerHandle()
	e.runStarter()
}

// Run Start Listen and serve
func (e *Engine) Run(addr string) error {
	e.init()
	return http.ListenAndServe(addr, e)
}

// GetMethodTree Getter for e.BluePrint.methodsTree
func (e *Engine) GetMethodTree() map[string][]*handleNode {
	return e.BluePrint.methodsTree
}

// ServeHTTP implement http.Handle
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := e.pool.Get().(*Context)
	group, params := e.Router.Match(request)
	if group != nil {
		context.matched = true
		group = append(e.Interceptors, group...)
	} else {
		group = append(e.Interceptors, e.NotFoundHandle)
	}
	context.init(request, writer, params, group)
	context.start()
	context.reset()
	e.pool.Put(context)
}

// New Constructor for Engine
func New() *Engine {
	engine := &Engine{
		Router:          make(HttpRouter),
		BluePrint:       NewBranch(),
		NotFoundHandle:  HandleNotFound,
		Warehouse:       new(SyncMap),
		MultipartMemory: defaultMultipartMemory,
		Abort:           exit{},
	}
	engine.pool = sync.Pool{New: func() interface{} { return engine.dispatchContext() }}
	engine.Use(HandleWithParser(JsonParser{}, FormParser{}, MultipartFormParser{}))
	return engine
}

// Default Engine for use
func Default() *Engine {
	engine := New()
	engine.AddInterceptors(LogInterceptor)
	engine.AddStarter(&BannerStarter{Banner: Banner}, &UrlInfoStarter{})
	return engine
}

// Map is a shortcut fot map[string]interface{}
type Map map[string]interface{}

func stringToByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&*(*reflect.StringHeader)(unsafe.Pointer(&s))))
}
