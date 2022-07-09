// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"net/http"
)

type HandleFunc func(context *Context)

type HandleFuncGroup []HandleFunc

func HandleNotFound(context *Context) { http.NotFound(context.ResponseWriter, context.Request) }

func RawHandlerFunc(handler http.HandlerFunc) HandleFunc {
	return func(c *Context) {
		SetContextIntoRequest(c)
		handler(c.ResponseWriter, c.Request)
	}
}

func RawHandlerFuncGroup(handlers ...http.HandlerFunc) HandleFuncGroup {
	group := make(HandleFuncGroup, len(handlers))
	for index, handler := range handlers {
		h := RawHandlerFunc(handler)
		group[index] = h
	}
	return group
}

// RecoverHandler recovers from panics and call given handler
func RecoverHandler(h func(ctx *Context, rec interface{})) HandleFunc {
	return func(context *Context) {
		defer func() {
			if rec := recover(); rec != nil {
				h(context, rec)
			}
		}()
		context.Next()
	}
}
