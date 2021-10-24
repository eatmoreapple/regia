// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"net/http"
)

type Router interface {
	Insert(method, path string, handle HandleFuncGroup)
	Match(req *http.Request) (HandleFuncGroup, Params)
}

// HttpRouter implement Router
type HttpRouter map[string]*routerNode

func (r HttpRouter) Insert(method, path string, handle HandleFuncGroup) {

	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	root := r[method]
	if root == nil {
		root = new(routerNode)
		r[method] = root
	}

	root.addRoute(path, handle)
}

func (r HttpRouter) Match(req *http.Request) (HandleFuncGroup, Params) {
	method := req.Method
	if root := r[method]; root != nil {
		group, params, _ := root.getValue(req.URL.Path)
		return group, params
	}
	return nil, nil
}
