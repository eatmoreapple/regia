// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"context"
	"errors"
	"github.com/eatmoreapple/regia/binders"
	"github.com/eatmoreapple/regia/logger"
	"github.com/eatmoreapple/regia/renders"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const defaultMultipartMemory = 32 << 20

type Context struct {
	// if url is matched
	matched bool
	// escape is a flag decide if is return to the context pool
	escape     bool
	index      uint8
	abortIndex uint8
	status     int
	written    bool
	// query cache
	queryCache url.Values
	// form cache
	formCache      url.Values
	items          map[string]interface{}
	lock           sync.RWMutex
	group          handleFuncNodeGroup
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Engine         *Engine
	Params         Params
	fullPath       string
}

// reset current Context
func (c *Context) reset() {
	c.index = 0
	c.matched = false
	c.queryCache = nil
	c.formCache = nil
	c.status = 0
	c.written = false
	c.abortIndex = 0
}

// start to handle current request
func (c *Context) start() {
	defer c.finish()
	c.Next()
}

// I do not think it is a good design
func (c *Context) finish() {
	if c.status != 0 && !c.written {
		c.ResponseWriter.WriteHeader(c.status)
	}
}

// IsMatched return that route matched
func (c *Context) IsMatched() bool {
	return c.matched
}

// Next call handle
func (c *Context) Next() {
	c.index++
	for c.index <= uint8(len(c.group)) && !c.IsAborted() {
		handle := c.group[c.index-1]
		handle.HandleFunc(c)
		c.index++
	}
}

// Flusher Make http.ResponseWriter as http.Flusher
func (c *Context) Flusher() http.Flusher { return c.ResponseWriter.(http.Flusher) }

// SaveUploadFile will call Context.FileStorage
// default save file to local path
func (c *Context) SaveUploadFile(name string) (string, error) {
	fs := c.BluePrint().FileStorage()
	return c.SaveUploadFileWith(fs, name)
}

// SaveUploadFileWith call given FileStorage with upload file
func (c *Context) SaveUploadFileWith(fs FileStorage, name string) (string, error) {
	if fs == nil {
		return "", errors.New("`FileStorage` can be nil type")
	}
	file, fileHeader, err := c.Request.FormFile(name)
	if err != nil {
		return "", err
	}
	// try to close return file
	if err = file.Close(); err != nil {
		return "", err
	}
	return fs.Save(fileHeader)
}

// Data analysis request body to destination and validate
// Call Context.AddParser to add more support
func (c *Context) Data(v interface{}) error {
	if err := c.BluePrint().Parsers().Parse(c, v); err != nil {
		return err
	}
	return c.BluePrint().Validator().Validate(v)
}

// Query is a shortcut for c.Request.URL.Query()
// but cached value for current context
func (c *Context) Query() url.Values {
	if c.queryCache == nil {
		c.queryCache = c.Request.URL.Query()
	}
	return c.queryCache
}

// QueryValue get Value from url query
func (c *Context) QueryValue(key string) Value {
	value := c.Query().Get(key)
	return Value(value)
}

// QueryValues get Value slice from url query
func (c *Context) QueryValues(key string) Values {
	values := c.Query()[key]
	return NewValues(values)
}

// Form is a shortcut for c.Request.PostForm
// but value for current context
func (c *Context) Form() url.Values {
	if c.formCache == nil {
		_ = c.Request.ParseForm()
		c.formCache = c.Request.PostForm
	}
	return c.formCache
}

// FormValue get Value from post value
func (c *Context) FormValue(key string) Value {
	value := c.Form().Get(key)
	return Value(value)
}

// FormValues get Values slice from post value
func (c *Context) FormValues(key string) Values {
	value := c.Form()[key]
	return NewValues(value)
}

// FullPath return full path of current request
func (c *Context) FullPath() string {
	return c.fullPath
}

// ContentType return Content-Type header
func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
}

// Bind bind request to destination
func (c *Context) Bind(binder binders.Binder, v interface{}) error {
	return binder.Bind(c.Request, v)
}

// BindQuery bind Query to destination
func (c *Context) BindQuery(v interface{}) error {
	binder := binders.QueryBinder{}
	return c.Bind(binder, v)
}

// BindForm bind PostForm to destination
func (c *Context) BindForm(v interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}
	binder := binders.FormBinder{}
	return c.Bind(binder, v)
}

// BindMultipartForm bind MultipartForm to destination
func (c *Context) BindMultipartForm(v interface{}) error {
	if err := c.Request.ParseMultipartForm(c.Engine.MultipartMemory); err != nil {
		return err
	}
	binder := binders.MultipartFormBodyBinder{}
	return c.Bind(binder, v)
}

// BindJSON bind the request body according to the format of json
func (c *Context) BindJSON(v interface{}) error {
	serializer := c.BluePrint().JSONSerializer()
	binder := binders.JsonBodyBinder{Serializer: serializer}
	return c.Bind(binder, v)
}

// BindXML bind the request body according to the format of xml
func (c *Context) BindXML(v interface{}) error {
	serializer := c.BluePrint().XMLSerializer()
	binder := binders.XmlBodyBinder{Serializer: serializer}
	return c.Bind(binder, v)
}

// BindHeader bind the request header to destination
func (c *Context) BindHeader(v interface{}) error {
	binder := binders.HeaderBinder{}
	return c.Bind(binder, v)
}

// BindURI bind the request uri to destination
func (c *Context) BindURI(v interface{}) error {
	values := c.Params.ToURLValues()
	binder := binders.URIBinder{Values: values}
	return c.Bind(binder, v)
}

// GetValue get value from context
func (c *Context) GetValue(key string) (value interface{}, exist bool) {
	c.lock.RLock()
	value, exist = c.items[key]
	c.lock.RUnlock()
	return
}

// SetValue set value to context
func (c *Context) SetValue(key string, value interface{}) {
	c.lock.Lock()
	if c.items == nil {
		c.items = make(map[string]interface{})
	}
	c.items[key] = value
	c.lock.Unlock()
}

// SetStatus set response status code
func (c *Context) SetStatus(code int) {
	if code < 0 {
		code = 0
	}
	c.status = code
}

// SetHeader set response header
func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// SetCookie is a shortcut for http.SetCookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

//************************
//*** Response Renders ***
//************************

// Render write response data with given Render
func (c *Context) Render(render renders.Render, data interface{}) error {
	if c.status != 0 {
		c.ResponseWriter.WriteHeader(c.status)
		c.status = 0
	}
	if !bodyAllowedForStatus(c.status) {
		return nil
	}
	c.written = true
	return render.Render(c.ResponseWriter, data)
}

// JSON write json response
func (c *Context) JSON(data interface{}) error {
	serializer := c.BluePrint().JSONSerializer()
	render := renders.JsonRender{Serializer: serializer}
	return c.Render(render, data)
}

// XML write xml response
func (c *Context) XML(data interface{}) error {
	serializer := c.BluePrint().XMLSerializer()
	render := renders.XmlRender{Serializer: serializer}
	return c.Render(render, data)
}

// String write string response
func (c *Context) String(format string, data ...interface{}) (err error) {
	render := renders.StringRender{Format: format, Data: data}
	return c.Render(render, nil)
}

// HTML write html response
func (c *Context) HTML(name string, data interface{}) error {
	render, err := c.BluePrint().HTMLLoader().Load(name)
	if err != nil {
		return err
	}
	return c.Render(render, data)
}

// Redirect Shortcut for http.Redirect
func (c *Context) Redirect(code int, url string) error {
	render := renders.RedirectRender{Code: code, RedirectURL: url, Request: c.Request}
	return c.Render(render, nil)
}

// ServeFile Shortcut for http.ServeFile
func (c *Context) ServeFile(filepath, filename string) error {
	render := renders.FileAttachmentRender{Request: c.Request, Filename: filename, FilePath: filepath}
	return c.Render(render, nil)
}

// ServeContent Shortcut for http.ServeContent
func (c *Context) ServeContent(name string, modTime time.Time, content io.ReadSeeker) error {
	render := renders.ContentRender{Name: name, ModTime: modTime, Content: content}
	return c.Render(render, nil)
}

// Write []byte into response writer
func (c *Context) Write(data []byte) error {
	_, err := c.ResponseWriter.Write(data)
	return err
}

// Escape can let context not return to the pool
func (c *Context) Escape() {
	c.escape = true
}

// IsEscape return escape status
func (c *Context) IsEscape() bool {
	return c.escape
}

//************************
//*** Abort Methods ******
//************************

// Abort skip current handle and will call Context.abort
// exit and do nothing by default
func (c *Context) Abort() { c.abortIndex = c.index }

// IsAborted return that context is aborted
func (c *Context) IsAborted() bool {
	return c.abortIndex != 0
}

// AbortHandler returns a handler which called at lasted
func (c *Context) AbortHandler() HandleFunc {
	if !c.IsAborted() {
		return nil
	}
	return c.group[c.abortIndex].HandleFunc
}

// AbortWithJSON write json response and exit
func (c *Context) AbortWithJSON(data interface{}) {
	_ = c.JSON(data)
	c.Abort()
}

// AbortWithXML write xml response and exit
func (c *Context) AbortWithXML(data interface{}) {
	_ = c.XML(data)
	c.Abort()
}

// AbortWithString write string response and exit
func (c *Context) AbortWithString(text string, data ...interface{}) {
	_ = c.String(text, data...)
	c.Abort()
}

// AbortWithStatus set response status and exit
func (c *Context) AbortWithStatus(code int) {
	c.SetStatus(code)
	c.Abort()
}

// IsWebsocket returns true if the request headers indicate that a websocket
func (c *Context) IsWebsocket() bool {
	return strings.Contains(strings.ToLower(c.Request.Header.Get("Connection")), "upgrade") &&
		strings.EqualFold(c.Request.Header.Get("Upgrade"), "websocket")
}

// IsAjax check current if is an ajax request
func (c *Context) IsAjax() bool {
	return strings.EqualFold(c.Request.Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// BluePrint return current blueprint
func (c *Context) BluePrint() *BluePrint {
	return c.group[c.index-1].BluePrint
}

func (c *Context) Logger() logger.Logger {
	return c.BluePrint().Logger()
}

type contextKey struct{}

type contextExistKey struct{}

// ContextKey is the request context key under which Context are stored.
var (
	ContextKey   = contextKey{}
	contextExist = contextExistKey{}
)

// SetContextIntoRequest set Context into request context
func SetContextIntoRequest(ctx *Context) {
	c := ctx.Request.Context()
	// if is the first time
	// ensure that called only once
	if c.Value(contextExist) == nil {
		c = context.WithValue(c, contextExist, contextExist)
		c = context.WithValue(c, ContextKey, ctx)
		ctx.Request = ctx.Request.WithContext(c)
	}
}

// GetCurrentContext get current Context from the request
func GetCurrentContext(req *http.Request) *Context {
	p, _ := req.Context().Value(ContextKey).(*Context)
	return p
}
