// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"github.com/eatmoreapple/regia/logger"
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

	// All requests will be intercepted by interceptors
	// whatever route matched or not
	interceptors HandleFuncGroup

	// Starter will run when the service starts
	// it only runs once
	starters []Starter

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
	HTMLLoader HTMLLoader

	// Logger used to log
	Logger logger.Logger

	// http.Server instance
	server *http.Server
}

func (e *Engine) dispatchContext() *Context {
	return &Context{
		Engine: e,
	}
}

// AddInterceptors Add interceptor to Engine
// All interceptors will be called before any handler
// Such as authorization, rate limiter, etc
func (e *Engine) AddInterceptors(interceptors ...HandleFunc) {
	e.interceptors = append(e.interceptors, interceptors...)
}

// AddStarter Add starter to Engine
// It will be called when the service starts
func (e *Engine) AddStarter(starters ...Starter) {
	e.starters = append(e.starters, starters...)
}

// init engine
func (e *Engine) init() error {
	// prepare router
	for method, nodes := range e.methodsTree {
		for _, node := range nodes {
			e.Router.Insert(method, node.path, node.group)
		}
	}
	// run all starters
	for _, starter := range e.starters {
		if err := starter.Start(e); err != nil {
			return err
		}
	}
	return nil
}

// Run is a shortcut for ListenAndServe
func (e *Engine) Run(addr string) error {
	return e.ListenAndServe(addr)
}

// ServeHTTP implement http.Handle
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := e.pool.Get().(*Context)
	context.Request = request
	context.ResponseWriter = writer

	// try to find all handlers
	group, params := e.Router.Match(context)

	context.matched = group != nil

	// if matched, then call the handler
	if context.matched {
		if len(e.interceptors) != 0 {
			group = append(e.interceptors, group...)
		}
	} else {
		// route not found
		// add not found handler
		// in case of not found handler is not set
		// then reply with 404
		// try to set Engine.NotFoundHandle to do your own business
		group = []HandleFunc{e.NotFoundHandle}
	}

	// initialize context
	context.init(params, group)

	// start to call all handlers
	context.start()

	// release context
	if !context.escape {
		context.reset()
		e.pool.Put(context)
	}
}

// ListenAndServeTLS acts identically to Run
func (e *Engine) ListenAndServeTLS(addr, certFile, keyFile string) error {
	if err := e.setup(); err != nil {
		return err
	}
	e.server.Addr = addr
	return e.server.ListenAndServeTLS(certFile, keyFile)
}

func (e *Engine) ListenAndServe(addr string) error {
	if err := e.setup(); err != nil {
		return err
	}
	e.server.Addr = addr
	return e.server.ListenAndServe()
}

// Server is a getter for Engine
func (e *Engine) Server() *http.Server {
	return e.server
}

func (e *Engine) CloneServer() *http.Server {
	return &http.Server{Handler: e}
}

func (e *Engine) setup() error {
	if err := e.init(); err != nil {
		return err
	}
	e.server = &http.Server{Handler: e}
	return nil
}

// New Constructor for Engine
func New() *Engine {
	engine := &Engine{
		Router:           HttpRouter{},
		BluePrint:        NewBluePrint(),
		NotFoundHandle:   HandleNotFound,
		Warehouse:        warehouse{},
		MultipartMemory:  defaultMultipartMemory,
		FileStorage:      &LocalFileStorage{},
		ContextValidator: validators.DefaultValidator{},
		HTMLLoader:       &TemplateLoader{},
		// Add default parser to make sure that Context could be worked
		ContextParser: Parsers{JsonParser{}, FormParser{}, MultipartFormParser{}, XMLParser{}},
		Logger:        logger.ConsoleLogger(),
	}
	engine.pool = sync.Pool{New: func() interface{} { return engine.dispatchContext() }}
	return engine
}

// Default Engine for use
func Default() *Engine {
	engine := New()
	engine.AddStarter(&BannerStarter{Banner: Banner}, &UrlInfoStarter{})
	return engine
}

const (
	author = "多吃点苹果"
	wechat = "eatmoreapple"
)
