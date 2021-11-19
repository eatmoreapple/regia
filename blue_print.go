// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
)

type handleNode struct {
	path  string
	group HandleFuncGroup
}

type BluePrint struct {
	methodsTree map[string][]*handleNode
	middleware  HandleFuncGroup
	prefix      string
}

// Use add middleware for this BluePrint
func (b *BluePrint) Use(group ...HandleFunc) { b.middleware = append(b.middleware, group...) }

// SetPrefix add prefix for this BluePrint
func (b *BluePrint) SetPrefix(path string) { b.prefix = path }

// GET is a shortcut for Handle("GET", path, group...)
func (b *BluePrint) GET(path string, group ...HandleFunc) {
	b.Handle(http.MethodGet, path, group...)
}

// POST is a shortcut for Handle("POST", path, group...)
func (b *BluePrint) POST(path string, group ...HandleFunc) {
	b.Handle(http.MethodPost, path, group...)
}

// PUT is a shortcut for Handle("PUT", path, group...)
func (b *BluePrint) PUT(path string, group ...HandleFunc) {
	b.Handle(http.MethodPut, path, group...)
}

// PATCH is a shortcut for Handle("PATCH", path, group...)
func (b *BluePrint) PATCH(path string, group ...HandleFunc) {
	b.Handle(http.MethodPatch, path, group...)
}

// DELETE is a shortcut for Handle("DELETE", path, group...)
func (b *BluePrint) DELETE(path string, group ...HandleFunc) {
	b.Handle(http.MethodDelete, path, group...)
}

// HEAD is a shortcut for Handle("HEAD", path, group...)
func (b *BluePrint) HEAD(path string, group ...HandleFunc) {
	b.Handle(http.MethodHead, path, group...)
}

// OPTIONS is a shortcut for Handle("OPTIONS", path, group...)
func (b *BluePrint) OPTIONS(path string, group ...HandleFunc) {
	b.Handle(http.MethodOptions, path, group...)
}

// ANY register all method for given handle
func (b *BluePrint) ANY(path string, group ...HandleFunc) {
	for _, method := range httpMethods {
		b.Handle(method, path, group...)
	}
}

// RAW register http.HandlerFunc with all method
func (b *BluePrint) RAW(method, path string, handlers ...http.HandlerFunc) {
	h := RawHandlerFuncGroup(handlers...)
	b.Handle(method, path, h...)
}

// Handle register HandleFunc with given method and path
func (b *BluePrint) Handle(method, path string, group ...HandleFunc) {
	group = append(b.middleware, group...)
	path = b.prefix + path
	n := &handleNode{path: path, group: group}
	if b.methodsTree == nil {
		b.methodsTree = make(map[string][]*handleNode)
	}
	methods := []string{method}
	if method == ALLMethods {
		methods = httpMethods[:]
	}
	for _, m := range methods {
		m = strings.ToUpper(m)
		b.methodsTree[m] = append(b.methodsTree[m], n)
	}
}

// Include can add another BluePrint
func (b *BluePrint) Include(prefix string, branch *BluePrint) {
	for method, nodes := range branch.methodsTree {
		for _, node := range nodes {
			b.Handle(method, prefix+node.path, node.group...)
		}
	}
}

// Bind legal handleFunc with given mappings from struct
func (b *BluePrint) Bind(path string, v interface{}, mappings ...map[string]string) {
	for _, mapping := range mappings {
		cleanedMapping := getCleanedRequestMapping(mapping)
		value := reflect.ValueOf(v)
		for handleName, methodName := range cleanedMapping {
			if method := value.MethodByName(handleName); method.IsValid() {
				if handle, ok := method.Interface().(func(context *Context)); ok {
					b.Handle(methodName, path, handle)
				}
			}
		}
	}
}

// BindMethod add HttpRequestMethodMapping for bind mappings
func (b *BluePrint) BindMethod(path string, v interface{}, mappings ...map[string]string) {
	mappings = append(mappings, HttpRequestMethodMapping)
	b.Bind(path, v, mappings...)
}

// BindByHandlerName register handler by handler name
// 		type Handler struct{}
// 		func(Handler)PostLogin(c *Context) {}
//		engine.BindByHandlerName("/user/", Handler{})
func (b *BluePrint) BindByHandlerName(path string, v interface{}) {
	value := reflect.ValueOf(v)
	t := reflect.TypeOf(v)
	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		methodName := t.Method(i).Name
		if m, ok := method.Interface().(func(ctx *Context)); ok {
			for k, v := range HttpRequestMethodMapping {
				if strings.HasPrefix(methodName, k) {
					pathName := strings.TrimLeft(methodName, k)
					b.Handle(v, path+getHandlerPathName(pathName), m)
				}
			}
		}
	}
}

// Static Serve static files
//        BluePrint.Static("/static/", "./static")
func (b *BluePrint) Static(url, dir string, group ...HandleFunc) {
	if strings.Contains(url, "*") {
		panic("`url` should not have wildcards")
	}
	server := http.FileServer(http.Dir(dir))
	handle := func(context *Context) {
		path := context.Params.Get(FilePathParam).Text()
		context.Request.URL.Path = path
		p := filepath.Join(dir, path)
		if _, err := os.Stat(p); err != nil {
			context.matched = false
			context.Engine.NotFoundHandle(context)
			context.AbortWith(exit{})
			return
		}
		ext := filepath.Ext(path)
		cnt := mime.TypeByExtension(ext)
		if len(cnt) == 0 {
			cnt = octetStream
		}
		context.SetHeader(contentType, cnt)
		server.ServeHTTP(context.ResponseWriter, context.Request)
	}
	group = append(group, handle)
	if !strings.HasSuffix(url, FilePathParam) {
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
	}
	url += wildFilepath
	b.Handle(http.MethodGet, url, group...)
}

func NewBluePrint() *BluePrint {
	return &BluePrint{}
}

func getCleanedRequestMapping(mapping map[string]string) map[string]string {
	cleanedMapping := make(map[string]string)
	for handleName, requestMethod := range mapping {
		requestMethodUpper := strings.ToUpper(requestMethod)
		for index, method := range httpMethods {
			if requestMethodUpper == method {
				break
			} else if index == (len(httpMethods)-1) && requestMethodUpper != method {
				panic("invalid method" + requestMethod)
			}
		}
		cleanedMapping[handleName] = requestMethodUpper
	}
	return cleanedMapping
}

func getHandlerPathName(name string) string {
	var builder strings.Builder
	for index, n := range name {
		if index == 0 {
			builder.WriteString(strings.ToLower(string(n)))
			continue
		}
		if unicode.IsUpper(n) {
			builder.WriteString("-")
			builder.WriteString(strings.ToLower(string(n)))
			continue
		} else {
			builder.WriteRune(n)
		}
	}
	return builder.String()
}
