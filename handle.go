package regia

import (
	"fmt"
	"net/http"
	"time"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

var logTitle = formatColor("[REGIA LOG]", colorGreen)

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

	defer func() {
		endTime := time.Since(start)
		startTimeStr := formatColor(start.Format(defaultTimeFormat), colorYellow)
		method := formatColor(fmt.Sprintf("[METHOD:%s]", context.Request.Method), colorBlue)
		path := formatColor(fmt.Sprintf("[PATH:%s]", context.Request.URL.Path), 96) // #02F3F3
		addr := formatColor(fmt.Sprintf("[Addr:%s]", context.Request.RemoteAddr), 97)
		end := formatColor(endTime.String(), colorMagenta)
		// 2006-01-02 15:04:05     [METHOD:GET]     [Addr:127.0.0.1:49453]      [PATH:/name]
		fmt.Printf("%-20s %-32s %-20s %-28s %-35s  %-20s\n", logTitle, startTimeStr, end, method, addr, path)
	}()

	context.Next()
}
