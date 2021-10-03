// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/regia/internal"
	"github.com/eatmoreapple/regia/validators"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HandleFunc func(context *Context)

type HandleFuncGroup []HandleFunc

func HandleWithValue(key string, value interface{}) HandleFunc {
	return func(context *Context) { context.ContextValue().Set(key, value) }
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

const timeFormat = "2006-01-02 15:04:05"

func LogInterceptor(context *Context) {
	start := time.Now()
	defer func() {
		endTime := time.Since(start)
		startTimeStr := internal.YellowString(start.Format(timeFormat))
		matched := internal.FormatColor(100, fmt.Sprintf("[MATCHED:%s]", strconv.FormatBool(context.IsMatched())))
		method := internal.FormatColor(101, fmt.Sprintf("[METHOD:%s]", context.Request.Method))
		path := internal.MagentaString(fmt.Sprintf("[PATH:%s]", context.Request.URL.Path)) // #02F3F3
		host, _, _ := net.SplitHostPort(context.Request.RemoteAddr)
		addr := internal.BlueString(fmt.Sprintf("[Addr:%s]", host))
		end := internal.MagentaString(endTime.String())
		// 2006-01-02 15:04:05     [METHOD:GET]     [Addr:127.0.0.1:49453]      [PATH:/name]
		fmt.Printf("%-32s %-20s %-28s %-28s %-40s %s\n", startTimeStr, end, matched, method, path, addr)
	}()
	context.Next()
}
