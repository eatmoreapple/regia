package logger

import (
	"fmt"
	"github.com/eatmoreapple/regia/internal"
	"io"
	"log"
)

type Logger interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type sampleLogger struct {
	*log.Logger
}

func (s *sampleLogger) Println(a ...interface{}) {
	s.Output(3, fmt.Sprintln(a...))
}

type logger struct {
	trace *sampleLogger
	debug *sampleLogger
	info  *sampleLogger
	warn  *sampleLogger
	error *sampleLogger
}

func (c *logger) Trace(v ...interface{}) {
	c.trace.Println(v...)
}

func (c *logger) Debug(v ...interface{}) {
	c.debug.Println(v...)
}

func (c *logger) Info(v ...interface{}) {
	c.info.Println(v...)
}

func (c *logger) Warn(v ...interface{}) {
	c.warn.Println(v...)
}

func (c *logger) Error(v ...interface{}) {
	c.error.Println(v...)
}

func newStdLogger(prefix string) *sampleLogger {
	return newLogger(log.Writer(), prefix)
}

func newLogger(writer io.Writer, prefix string) *sampleLogger {
	return &sampleLogger{log.New(writer, fmt.Sprintf("[%-14s]  ", prefix), log.Ldate|log.Ltime|log.Llongfile)}
}

func ConsoleLogger() Logger {
	trace := newStdLogger(internal.FormatColor(97, "TRACE"))
	debug := newStdLogger(internal.FormatColor(91, "DEBUG"))
	info := newStdLogger(internal.FormatColor(92, "INFO"))
	warn := newStdLogger(internal.FormatColor(93, "WARN"))
	err := newStdLogger(internal.FormatColor(91, "ERROR"))
	return NewLogger(trace, debug, info, warn, err)
}

func NewLogger(trace, debug, info, warn, error *sampleLogger) Logger {
	return &logger{
		trace: trace,
		debug: debug,
		info:  info,
		warn:  warn,
		error: error,
	}
}
