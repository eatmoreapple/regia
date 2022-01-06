// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"github.com/eatmoreapple/regia/validators"
	"net/http"
	"sync"
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

	// global Context FileStorage
	FileStorage FileStorage

	// global Context ContextValidator
	ContextValidator validators.Validator

	// global ContextParser
	ContextParser Parsers

	// HTML Loader
	TemplateLoader TemplateLoader
	server         *http.Server
}

func (e *Engine) dispatchContext() *Context {
	return &Context{
		Engine: e,
	}
}

// Start implement Starter and register all handles to router
func (e *Engine) Start(*Engine) error {
	for method, nodes := range e.methodsTree {
		for _, node := range nodes {
			e.Router.Insert(method, node.path, node.group)
		}
	}
	return nil
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
func (e *Engine) runStarter() error {
	for _, starter := range e.Starters {
		if err := starter.Start(e); err != nil {
			return err
		}
	}
	return nil
}

// Init engine
func (e *Engine) init() error {
	e.AddStarter(e)
	return e.runStarter()
}

// Run Start Listen and serve
func (e *Engine) Run(addr string) error {
	return e.ListenAndServe(addr)
}

// ServeHTTP implement http.Handle
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := e.pool.Get().(*Context)
	group, params := e.Router.Match(request)
	context.matched = len(group) == 0
	if group != nil {
		if len(e.Interceptors) != 0 {
			group = append(e.Interceptors, group...)
		}
	} else {
		group = []HandleFunc{e.NotFoundHandle}
		if len(e.Interceptors) != 0 {
			group = append(e.Interceptors, group...)
		}
	}
	context.init(request, writer, params, group)
	context.start()
	if !context.escape {
		context.reset()
		e.pool.Put(context)
	}
}

// ListenAndServeTLS acts identically to Run
func (e *Engine) ListenAndServeTLS(addr, certFile, keyFile string) error {
	err := e.SetUp()
	if err != nil {
		return err
	}
	e.server.Addr = addr
	err = e.server.ListenAndServeTLS(certFile, keyFile)
	return err
}

func (e *Engine) ListenAndServe(addr string) error {
	err := e.SetUp()
	if err != nil {
		return err
	}
	e.server.Addr = addr
	err = e.server.ListenAndServe()
	return err
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

func (e *Engine) SetUp() error {
	if err := e.init(); err != nil {
		return err
	}
	e.makeServer()
	return nil
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
		FileStorage:               &LocalFileStorage{},
		ContextValidator:          validators.DefaultValidator{},
		TemplateLoader:            new(HTMLLoader),
		// Add default parser to make sure that Context could be worked
		ContextParser: Parsers{JsonParser{}, FormParser{}, MultipartFormParser{}, XMLParser{}},
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
