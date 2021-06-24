package regia

import (
	"fmt"
	"net/http"
	"strconv"
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

func regiaLog(start time.Time, context *Context) {
	endTime := time.Since(start)
	startTimeStr := formatColor(colorYellow, start.Format(defaultTimeFormat))
	matched := formatColor(100, fmt.Sprintf("[MATCHED:%s]", strconv.FormatBool(context.IsMatched())))
	method := formatColor(102, fmt.Sprintf("[METHOD:%s]", context.Request.Method))
	path := formatColor(101, fmt.Sprintf("[PATH:%s]", context.Request.URL.Path)) // #02F3F3
	addr := formatColor(104, fmt.Sprintf("[Addr:%s]", context.Request.RemoteAddr))
	end := formatColor(colorMagenta, endTime.String())
	// 2006-01-02 15:04:05     [METHOD:GET]     [Addr:127.0.0.1:49453]      [PATH:/name]
	fmt.Printf("%-23s %-32s %-20s %-28s %-28s %-35s %-20s\n", logTitle, startTimeStr, end, matched, method, addr, path)
}
