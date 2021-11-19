// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"github.com/eatmoreapple/regia/validators"
	"net/http"
	"reflect"
	"sync"
	"unsafe"
)

const (
	FilePathParam = "static"
	wildFilepath  = "*" + FilePathParam
)

// Engine is a collection of core components of the whole service
type Engine struct {
	// BluePrint is used for store the handler.
	// All handlers are going to register into it
	*BluePrint

	// Router is a module used to register handle and distribute request
	Router Router

	// NotFoundHandle replies to the request with an HTTP 404 not found error.
	NotFoundHandle func(context *Context)

	// InternalServerErrorHandle replies to the request with an HTTP 500 internal server error.
	InternalServerErrorHandle func(context *Context, rec interface{})

	// All requests will be intercepted by Interceptors
	// whatever route matched or not
	Interceptors HandleFuncGroup

	// Starter will run when the service starts
	// it only runs once
	Starters []Starter

	// Warehouse is used to store information
	Warehouse Warehouse

	// MultipartMemory defined max request body size
	MultipartMemory int64

	// Context pool
	pool sync.Pool

	// global Context Abort
	Abort Exit

	// global Context FileStorage
	FileStorage FileStorage

	// global Context ContextValidator
	ContextValidator validators.Validator

	// global ContextParser
	ContextParser Parsers
	server        *http.Server
}

func (e *Engine) dispatchContext() *Context {
	return &Context{
		Engine: e,
	}
}

// Start implement Starter and register all handles to router
func (e *Engine) Start(*Engine) {
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
	e.AddStarter(e)
	e.runStarter()
}

// Run Start Listen and serve
func (e *Engine) Run(addr string) error {
	return e.ListenAndServe(addr)
}

// ServeHTTP implement http.Handle
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := e.pool.Get().(*Context)
	group, params := e.Router.Match(request)
	if group != nil {
		context.matched = true
		if len(e.Interceptors) != 0 {
			group = append(e.Interceptors, group...)
		}
	} else {
		if len(e.Interceptors) != 0 {
			group = append(e.Interceptors, e.NotFoundHandle)
		}
	}
	context.init(request, writer, params, group)
	context.start()
	if !context.escape {
		context.reset()
		e.pool.Put(context)
	}
}

// New Constructor for Engine
func New() *Engine {
	engine := &Engine{
		Router:                    make(HttpRouter),
		BluePrint:                 NewBluePrint(),
		NotFoundHandle:            HandleNotFound,
		InternalServerErrorHandle: HandleInternalServerError,
		Warehouse:                 new(SyncMap),
		MultipartMemory:           defaultMultipartMemory,
		Abort:                     exit{},
		FileStorage:               &LocalFileStorage{},
		ContextValidator:          validators.DefaultValidator{},
		// Add default parser to make sure that Context could be worked
		ContextParser: Parsers{JsonParser{}, FormParser{}, MultipartFormParser{}},
	}
	engine.pool = sync.Pool{New: func() interface{} { return engine.dispatchContext() }}
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

// unsafe string to byte
// without memory copy
func stringToByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&*(*reflect.StringHeader)(unsafe.Pointer(&s))))
}

// ListenAndServeTLS acts identically to Run
func (e *Engine) ListenAndServeTLS(addr, certFile, keyFile string) error {
	e.SetUp()
	e.server.Addr = addr
	return e.server.ListenAndServeTLS(certFile, keyFile)
}

func (e *Engine) ListenAndServe(addr string) error {
	e.SetUp()
	e.server.Addr = addr
	return e.server.ListenAndServe()
}

// Server is a getter for Engine
func (e *Engine) Server() *http.Server {
	return e.server
}

func (e *Engine) makeServer() {
	e.server = e.CloneServer()
}

func (e *Engine) CloneServer() *http.Server {
	return &http.Server{Handler: e}
}

func (e *Engine) SetUp() {
	e.init()
	e.makeServer()
}
