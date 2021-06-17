package regia

import (
	"errors"
	"io"
	"net/http"
	"time"
)

const defaultMultipartMemory = 32 << 20

type Context struct {
	*http.Request
	http.ResponseWriter
	index int
	group HandleFuncGroup
	// Mat multipart form memory size
	// default 32M
	MultipartMemory int64
	contextValue    *SyncMap
	Engine          *Engine
	FileStorage     FileStorage
	Parsers         Parsers
	Authenticators  Authenticators
	Params          Params
	abort           Exit
}

func (c *Context) init(req *http.Request, writer http.ResponseWriter, params Params, group HandleFuncGroup) {
	c.Request = req
	c.ResponseWriter = writer
	c.Params = params
	c.group = group
	c.FileStorage = LocalFileStorage{}
	c.abort = c.Engine.Abort
	c.MultipartMemory = c.Engine.MultipartMemory
	c.index = 0
}

func (c *Context) reset() {
	c.index = 0
	c.Request = nil
	c.ResponseWriter = nil
	c.Params = nil
	c.group = nil
	c.FileStorage = nil
	c.abort = nil
	c.MultipartMemory = defaultMultipartMemory
	c.contextValue = nil
	c.Parsers = nil
	c.Authenticators = nil
}

func (c *Context) start() {
	defer c.recover()
	c.Next()
}

func (c *Context) recover() {
	if rec := recover(); rec != nil {
		if e, ok := rec.(Exit); ok {
			e.Exit(c)
		} else {
			panic(rec)
		}
	}
}

func (c *Context) Next() {
	c.index++
	for c.index <= len(c.group) {
		handle := c.group[c.index-1]
		handle(c)
		c.index++
	}
}

func (c *Context) SetAbort(abort Exit) {
	c.abort = abort
}

// Abort Just exit and do nothing
func (c *Context) Abort() { c.AbortWith(c.abort) }

func (c *Context) AbortWith(exit Exit) { panic(exit) }

// Flusher Make http.ResponseWriter as http.Flusher
func (c *Context) Flusher() http.Flusher { return c.ResponseWriter.(http.Flusher) }

// SaveUploadFile implement your own idea with it
func (c *Context) SaveUploadFile(name string) error {
	return c.SaveUploadFileWith(c.FileStorage, name)
}

func (c *Context) SaveUploadFileWith(fs FileStorage, name string) error {
	if fs == nil {
		return errors.New("`FileStorage` can be nil type")
	}
	file, fileHeader, err := c.Request.FormFile(name)
	if err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return c.FileStorage.Save(fileHeader)
}

func (c *Context) Data(v interface{}) error {
	return c.Parsers.Parse(c, v)
}

func (c *Context) AddParser(p ...Parser) {
	c.Parsers = append(c.Parsers, p...)
}

func (c *Context) User(v interface{}) error {
	return c.Authenticators.RunAuthenticate(c, v)
}

func (c *Context) AddAuthenticator(a ...Authenticator) {
	c.Authenticators = append(c.Authenticators, a...)
}

func (c *Context) ContextValue() *SyncMap {
	if c.contextValue == nil {
		c.contextValue = new(SyncMap)
	}
	return c.contextValue
}

func (c *Context) Bind(binder Binder, v interface{}) error {
	return binder.Bind(c, v)
}

func (c *Context) BindQuery(v interface{}) error {
	binder := QueryBinder{}
	return c.Bind(binder, v)
}

func (c *Context) BindForm(v interface{}) error {
	if err := c.ParseForm(); err != nil {
		return err
	}
	binder := FormBinder{}
	return c.Bind(binder, v)
}

func (c *Context) BindMultipartForm(v interface{}) error {
	if err := c.Request.ParseMultipartForm(c.MultipartMemory); err != nil {
		return err
	}
	binder := MultipartFormBinder{}
	return c.Bind(binder, v)
}

func (c *Context) BindJSON(v interface{}) error {
	binder := JsonBodyBinder{}
	return c.Bind(binder, v)
}

func (c *Context) BindXML(v interface{}) error {
	binder := XmlBodyBinder{}
	return c.Bind(binder, v)
}

func (c *Context) SetStatus(code int) {
	c.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.ResponseWriter, cookie)
}

func (c *Context) Render(render Render, data interface{}) error {
	return render.Render(c.ResponseWriter, data)
}

func (c *Context) JSON(data interface{}) error {
	render := JsonRender{}
	return c.Render(render, data)
}

func (c *Context) XML(data interface{}) error {
	render := XmlRender{}
	return c.Render(render, data)
}

func (c *Context) Text(text string) (int, error) {
	writeContentType(c.ResponseWriter, textHtmlContentType)
	return c.ResponseWriter.Write(stringToByte(text))
}

// Redirect Shortcut for http.Redirect
func (c *Context) Redirect(code int, url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, code)
}

// ServeFile Shortcut for http.ServeFile
func (c *Context) ServeFile(path string) {
	http.ServeFile(c.ResponseWriter, c.Request, path)
}

// ServeContent Shortcut for http.ServeContent
func (c *Context) ServeContent(name string, modTime time.Time, content io.ReadSeeker) {
	http.ServeContent(c.ResponseWriter, c.Request, name, modTime, content)
}
