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
	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		panic("`v` should be a struct")
	}
	for _, mapping := range mappings {
		cleanedMapping := getCleanedRequestMapping(mapping)
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
			context.Abort()
			return
		}
		ext := filepath.Ext(path)
		cnt := mime.TypeByExtension(ext)
		if len(cnt) == 0 {
			cnt = "application/octet-stream"
		}
		context.SetHeader("Content-Type", cnt)
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
