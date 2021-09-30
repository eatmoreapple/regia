// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

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

func formatColor(color int, v interface{}) string {
	return fmt.Sprintf(colorFormat, color, fmt.Sprintf("%+v", v))
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

type Logger interface {
	Info(v interface{})
	Debug(v interface{})
	Warn(v interface{})
	Error(v interface{})
}

// ColorLogger implement Log
type ColorLogger struct {
	*log.Logger
}

func (c ColorLogger) Info(v interface{}) {
	c.Println(formatColor(colorGreen, v))
}

func (c ColorLogger) Debug(v interface{}) {
	c.Println(formatColor(colorMagenta, v))
}

func (c ColorLogger) Warn(v interface{}) {
	c.Println(formatColor(colorYellow, v))
}

func (c ColorLogger) Error(v interface{}) {
	c.Println(formatColor(colorRed, v))
}

var Log Logger = ColorLogger{log.New(os.Stdout, "", 0)}
