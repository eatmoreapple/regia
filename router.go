// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

// HttpRouter implement Router
type HttpRouter map[string]*routerNode

func (r HttpRouter) Insert(method, path string, handle handleFuncNodeGroup) {

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

func (r HttpRouter) Match(ctx *Context) bool {
	method := ctx.Request.Method
	if root := r[method]; root != nil {
		group, params, _ := root.getValue(ctx.Request.URL.Path)
		ctx.fullPath = root.fullPath
		ctx.Params = params
		ctx.group = group
		return group != nil
	}
	return false
}
