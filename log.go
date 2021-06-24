package regia

import "fmt"

const colorFormat = "\u001B[%dm%s\u001B[0m"

const (
	colorRed     = iota + 91 // red
	colorGreen               // green
	colorYellow              // yellow
	colorBlue                // blue
	colorMagenta             // magenta
)

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
	fmt.Println(formatColor(colorGreen, text))
}

func (c ConsoleLog) Debug(text string) {
	fmt.Println(formatColor(colorMagenta, text))
}

func (c ConsoleLog) Warn(text string) {
	fmt.Println(formatColor(colorYellow, text))
}

func (c ConsoleLog) Error(text string) {
	fmt.Println(formatColor(colorRed, text))
}

var Logger Log = ConsoleLog{}
