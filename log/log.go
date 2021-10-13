// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/eatmoreapple/regia/internal"
	"log"
	"os"
)

type Logger interface {
	Info(v interface{})
	Debug(v interface{})
	Warn(v interface{})
	Error(v interface{})
}

func Info(v interface{}) {
	DefaultLogger.Info(v)
}

func Debug(v interface{}) {
	DefaultLogger.Debug(v)
}

func Warn(v interface{}) {
	DefaultLogger.Warn(v)
}

func Error(v interface{}) {
	DefaultLogger.Error(v)
}

var DefaultLogger Logger = ColorLogger{log.New(os.Stdout, "", 0)}

// ColorLogger implement DefaultLogger
type ColorLogger struct {
	*log.Logger
}

func (c ColorLogger) Info(v interface{}) {
	c.Println(internal.GreenString(v))
}

func (c ColorLogger) Debug(v interface{}) {
	c.Println(internal.MagentaString(v))
}

func (c ColorLogger) Warn(v interface{}) {
	c.Println(internal.YellowString(v))
}

func (c ColorLogger) Error(v interface{}) {
	c.Println(internal.RedString(v))
}
