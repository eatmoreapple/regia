package regia

import (
	"fmt"
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

type Log interface {
	Info(text string)
	Debug(text string)
	Warn(text string)
	Error(text string)
}

type ConsoleLog struct{}

func (c ConsoleLog) Info(text string) {
	text = c.Format(info, text)
	fmt.Println(formatColor(colorGreen, text))
}

func (c ConsoleLog) Debug(text string) {
	text = c.Format(debug, text)
	fmt.Println(formatColor(colorMagenta, text))
}

func (c ConsoleLog) Warn(text string) {
	text = c.Format(warn, text)
	fmt.Println(formatColor(colorYellow, text))
}

func (c ConsoleLog) Error(text string) {
	text = c.Format(danger, text)
	fmt.Println(formatColor(colorRed, text))
}

func (c ConsoleLog) Format(level, text string) string {
	return fmt.Sprintf("[%-5s]  %-20s    %s", level, time.Now().Format(timeFormat), text)
}

var Logger Log = ConsoleLog{}
