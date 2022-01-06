// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"github.com/eatmoreapple/regia/validators"
	"net/http"
)

type HandleFunc func(context *Context)

type HandleFuncGroup []HandleFunc

func HandleWithValue(key string, value interface{}) HandleFunc {
	return func(context *Context) { context.SetValue(key, value) }
}

func HandleWithParser(parser ...Parser) HandleFunc {
	return func(context *Context) { context.AddParser(parser...) }
}

func HandleWithValidator(validator validators.Validator) HandleFunc {
	return func(context *Context) { context.Validator = validator }
}

func HandleWithFileStorage(fs FileStorage) HandleFunc {
	return func(context *Context) { context.FileStorage = fs }
}

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
