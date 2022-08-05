// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"github.com/eatmoreapple/regia/internal"
	"github.com/eatmoreapple/regia/logger"
	"github.com/eatmoreapple/regia/validators"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type handleNode struct {
	path      string
	group     HandleFuncGroup
	blueprint *BluePrint
}

type BluePrint struct {
	// Name the name of current BluePrint
	Name string

	// data of current BluePrint
	Data interface{}

	// FileStorage is a storage for file
	fileStorage FileStorage
	parsers     Parsers
	validator   validators.Validator
	logger      logger.Logger

	// response render
	htmlLoader    HTMLLoader
	xmlSerializer internal.Serializer
	// JSONSerializer used to serialize json
	// Your set your own JSONSerializer if you want
	// Such as jsoniter, json2, etc
	jsonSerializer internal.Serializer

	parent      *BluePrint
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
	n := &handleNode{path: path, group: group, blueprint: b}
	b.register(method, n)
}

// Register register handleNode with given method
func (b *BluePrint) register(method string, node *handleNode) {
	if b.methodsTree == nil {
		b.methodsTree = make(map[string][]*handleNode)
	}
	methods := []string{method}
	if method == ALLMethods {
		methods = httpMethods[:]
	}
	for _, m := range methods {
		m = strings.ToUpper(m)
		b.methodsTree[m] = append(b.methodsTree[m], node)
	}
}

// Include can add another BluePrint
func (b *BluePrint) Include(prefix string, branch *BluePrint) {
	// set parent
	branch.parent = b
	for method, nodes := range branch.methodsTree {
		for _, node := range nodes {
			hn := &handleNode{path: prefix + node.path, group: node.group, blueprint: branch}
			b.register(method, hn)
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
		// get file path from url
		path := context.Params.Get(FilePathParam).Text()
		// reset file path to URL
		context.Request.URL.Path = path

		// if file exists, serve it
		p := filepath.Join(dir, path)
		if _, err := os.Stat(p); err != nil {
			context.matched = false
			context.Engine.NotFoundHandle(context)
			context.Abort()
			return
		}
		ext := filepath.Ext(path)

		// try to serve file with correct content type
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

// Parent returns parent BluePrint
func (b *BluePrint) Parent() *BluePrint {
	return b.parent
}

// IsRoot returns true if this is root BluePrint
func (b *BluePrint) IsRoot() bool {
	return b.Parent() == nil
}

//*************************
//*** Getter And Setter ***
//*************************

// XMLSerializer returns XMLSerializer
// If not set, it will try to get from parent BluePrint
func (b *BluePrint) XMLSerializer() internal.Serializer {
	if b.xmlSerializer != nil {
		return b.xmlSerializer
	}
	if !b.IsRoot() {
		return b.Parent().XMLSerializer()
	}
	return nil
}

// SetXMLSerializer set XMLSerializer
// If is nil, it will be panic
func (b *BluePrint) SetXMLSerializer(xmlSerializer internal.Serializer) {
	if xmlSerializer == nil {
		panic("xmlSerializer can not be nil")
	}
	b.xmlSerializer = xmlSerializer
}

// JSONSerializer returns JSONSerializer
// If not set, it will try to get from parent BluePrint
func (b *BluePrint) JSONSerializer() internal.Serializer {
	if b.jsonSerializer != nil {
		return b.jsonSerializer
	}
	if !b.IsRoot() {
		return b.Parent().JSONSerializer()
	}
	return nil
}

// SetJSONSerializer set JSONSerializer
// If is nil, it will be panic
func (b *BluePrint) SetJSONSerializer(jsonSerializer internal.Serializer) {
	if jsonSerializer == nil {
		panic("jsonSerializer can not be nil")
	}
	b.jsonSerializer = jsonSerializer
}

// HTMLLoader HTMLSerializer returns HTMLLoader
// If not set, it will try to get from parent BluePrint
func (b *BluePrint) HTMLLoader() HTMLLoader {
	if b.htmlLoader != nil {
		return b.htmlLoader
	}
	if !b.IsRoot() {
		return b.Parent().HTMLLoader()
	}
	return nil
}

// SetHTMLLoader set HTMLLoader
// If is nil, it will be panic
func (b *BluePrint) SetHTMLLoader(htmlLoader HTMLLoader) {
	if htmlLoader == nil {
		panic("htmlLoader can not be nil")
	}
	b.htmlLoader = htmlLoader
}

// FileStorage set FileStorage
// If is nil, it will be panic
func (b *BluePrint) FileStorage() FileStorage {
	if b.fileStorage != nil {
		return b.fileStorage
	}
	if !b.IsRoot() {
		return b.Parent().FileStorage()
	}
	return nil
}

// SetFileStorage set FileStorage
// If is nil, it will be panic
func (b *BluePrint) SetFileStorage(fileStorage FileStorage) {
	if fileStorage == nil {
		panic("fileStorage can not be nil")
	}
	b.fileStorage = fileStorage
}

// Parsers returns Parsers
func (b *BluePrint) Parsers() Parsers {
	if b.parsers != nil {
		return b.parsers
	}
	if !b.IsRoot() {
		return b.Parent().Parsers()
	}
	return nil
}

// SetParsers set Parsers
// If is nil, it will be panic
func (b *BluePrint) SetParsers(parsers Parsers) {
	if parsers == nil {
		panic("parsers can not be nil")
	}
	b.parsers = parsers
}

// Validator returns Validator
// If not set, it will try to get from parent BluePrint
func (b *BluePrint) Validator() validators.Validator {
	if b.validator != nil {
		return b.validator
	}
	if !b.IsRoot() {
		return b.Parent().Validator()
	}
	return nil
}

// SetValidator set Validator
// If is nil, it will be panic
func (b *BluePrint) SetValidator(validator validators.Validator) {
	if validator == nil {
		panic("validator can not be nil")
	}
	b.validator = validator
}

// Logger returns Logger
// If not set, it will try to get from parent BluePrint
func (b *BluePrint) Logger() logger.Logger {
	if b.logger != nil {
		return b.logger
	}
	if !b.IsRoot() {
		return b.Parent().Logger()
	}
	return nil
}

// SetLogger set Logger
// If is nil, it will be panic
func (b *BluePrint) SetLogger(logger logger.Logger) {
	if logger == nil {
		panic("logger can not be nil")
	}
	b.logger = logger
}

// NewBluePrint constructor for BluePrint
func NewBluePrint() *BluePrint {
	return &BluePrint{}
}

// DefaultBluePrint returns default BluePrint
// It is used to create root BluePrint
// If you want to create child BluePrint, you should use NewBluePrint
// DefaultBluePrint add some default attributes to BluePrint
// Ensure that system could be worked
func DefaultBluePrint() *BluePrint {
	bp := NewBluePrint()
	bp.SetFileStorage(&LocalFileStorage{})
	bp.SetValidator(&validators.DefaultValidator{})
	bp.SetParsers(Parsers{JsonParser{}, FormParser{}, MultipartFormParser{}, XMLParser{}})
	bp.SetHTMLLoader(&TemplateLoader{})
	bp.SetLogger(logger.ConsoleLogger())
	bp.SetJSONSerializer(internal.JsonSerializer{})
	bp.SetXMLSerializer(internal.XmlSerializer{})
	return bp
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
