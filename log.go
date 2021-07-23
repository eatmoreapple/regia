package regia

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const colorFormat = "\u001B[%dm%s\u001B[0m"

const (
	colorRed     = iota + 91 // red
	colorGreen               // green
	colorYellow              // yellow
	colorBlue                // blue
	colorMagenta             // magenta
)

const (
	info   = "INFO"
	debug  = "DEBUG"
	warn   = "WARN"
	danger = "DANGER"
)

const timeFormat = "2006-01-02 15:04:05"

func formatColor(color int, text string) string {
	return fmt.Sprintf(colorFormat, color, text)
}

func regiaLog(start time.Time, context *Context) {
	endTime := time.Since(start)
	startTimeStr := formatColor(colorYellow, start.Format(defaultTimeFormat))
	matched := formatColor(100, fmt.Sprintf("[MATCHED:%s]", strconv.FormatBool(context.IsMatched())))
	method := formatColor(101, fmt.Sprintf("[METHOD:%s]", context.Request.Method))
	path := formatColor(colorMagenta, fmt.Sprintf("[PATH:%s]", context.Request.URL.Path)) // #02F3F3
	host, _, _ := net.SplitHostPort(context.Request.RemoteAddr)
	addr := formatColor(colorBlue, fmt.Sprintf("[Addr:%s]", host))
	end := formatColor(colorMagenta, endTime.String())
	// 2006-01-02 15:04:05     [METHOD:GET]     [Addr:127.0.0.1:49453]      [PATH:/name]
	fmt.Printf("%-23s %-32s %-20s %-28s %-28s %-40s %s\n", logTitle, startTimeStr, end, matched, method, path, addr)
}

type Log interface {
	Info(text string)
	Debug(text string)
	Warn(text string)
	Error(text string)
}

// ConsoleLog implement Log
type ConsoleLog struct {
	*log.Logger
}

func (c ConsoleLog) Info(text string) {
	c.Println(formatColor(colorGreen, text))
}

func (c ConsoleLog) Debug(text string) {
	c.Println(formatColor(colorMagenta, text))
}

func (c ConsoleLog) Warn(text string) {
	c.Println(formatColor(colorYellow, text))
}

func (c ConsoleLog) Error(text string) {
	c.Println(formatColor(colorRed, text))
}

var Logger Log = ConsoleLog{log.New(os.Stdout, "", 0)}
