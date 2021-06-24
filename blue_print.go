package regia

import (
	"net/http"
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

func (b *BluePrint) Use(group ...HandleFunc) { b.middleware = append(b.middleware, group...) }

func (b *BluePrint) SetPrefix(path string) { b.prefix = path }

func (b *BluePrint) GET(path string, group ...HandleFunc) {
	b.Handle(http.MethodGet, path, group...)
}

func (b *BluePrint) POST(path string, group ...HandleFunc) {
	b.Handle(http.MethodPost, path, group...)
}

func (b *BluePrint) PUT(path string, group ...HandleFunc) {
	b.Handle(http.MethodPut, path, group...)
}

func (b *BluePrint) PATCH(path string, group ...HandleFunc) {
	b.Handle(http.MethodPatch, path, group...)
}

func (b *BluePrint) DELETE(path string, group ...HandleFunc) {
	b.Handle(http.MethodDelete, path, group...)
}

func (b *BluePrint) HEAD(path string, group ...HandleFunc) {
	b.Handle(http.MethodHead, path, group...)
}

func (b *BluePrint) OPTIONS(path string, group ...HandleFunc) {
	b.Handle(http.MethodOptions, path, group...)
}

func (b *BluePrint) ANY(path string, group ...HandleFunc) {
	for _, method := range httpMethods {
		b.Handle(method, path, group...)
	}
}

func (b *BluePrint) Handle(method, path string, group ...HandleFunc) {
	group = append(b.middleware, group...)
	path = b.prefix + path
	n := &handleNode{path: path, group: group}
	if b.methodsTree == nil {
		b.methodsTree = make(map[string][]*handleNode)
	}
	b.methodsTree[method] = append(b.methodsTree[method], n)
}

func (b *BluePrint) Include(prefix string, branch *BluePrint) {
	for method, nodes := range branch.methodsTree {
		for _, node := range nodes {
			b.Handle(method, prefix+node.path, node.group...)
		}
	}
}

func (b *BluePrint) Bind(path string, v interface{}, mappings ...map[string]string) {
	for _, mapping := range mappings {
		cleanedMapping := b.getCleanedRequestMapping(mapping)
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

func (b *BluePrint) BindMethod(path string, v interface{}, mappings ...map[string]string) {
	mappings = append(mappings, HttpRequestMethodMapping)
	b.Bind(path, v, mappings...)
}

func (b *BluePrint) getCleanedRequestMapping(mapping map[string]string) map[string]string {
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

func NewBranch() *BluePrint {
	return &BluePrint{}
}
