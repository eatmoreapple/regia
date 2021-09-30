// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

var logTitle = formatColor(colorGreen, "[REGIA LOG]")

type HandleFunc func(context *Context)

type HandleFuncGroup []HandleFunc

func HandleWithValue(key string, value interface{}) HandleFunc {
	return func(context *Context) { context.ContextValue().Set(key, value) }
}

func HandleWithParser(parser ...Parser) HandleFunc {
	return func(context *Context) { context.AddParser(parser...) }
}

func HandleWithValidator(validator Validator) HandleFunc {
	return func(context *Context) { context.Validator = validator }
}

func HandleWithFileStorage(fs FileStorage) HandleFunc {
	return func(context *Context) { context.FileStorage = fs }
}

func HandleWithAbort(e Exit) HandleFunc {
	return func(context *Context) { context.abort = e }
}

func HandleOptions(allowMethods ...string) HandleFunc {
	if allowMethods == nil {
		allowMethods = httpMethods[:]
	}
	return func(context *Context) {
		if context.Request.Method == http.MethodOptions {
			context.SetHeader("Allow", strings.Join(allowMethods, ", "))
			context.SetHeader("Content-Length", "0")
			context.SetStatus(http.StatusNoContent)
			context.Abort()
		}
		context.Next()
	}
}

func HandleNotFound(context *Context) { http.NotFound(context.ResponseWriter, context.Request) }

func LogInterceptor(context *Context) {
	start := time.Now()

	defer regiaLog(start, context)

	context.Next()
}

func HandleWithExit(exit Exit) HandleFunc {
	return func(context *Context) {
		context.AbortWith(exit)
	}
}

func Flush(exit Exit) error {
	if exit == nil {
		return errors.New("exit can not be nil")
	}
	panic(exit)
}
