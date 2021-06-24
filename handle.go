package regia

import (
	"net/http"
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

func HandleNotFound(context *Context) { http.NotFound(context.ResponseWriter, context.Request) }

func LogInterceptor(context *Context) {
	start := time.Now()

	defer regiaLog(start, context)

	context.Next()
}


