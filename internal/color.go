package internal

import "fmt"

const colorFormat = "\u001B[%dm%s\u001B[0m"

const (
	Red     = iota + 91 // red
	Green               // green
	Yellow              // yellow
	Blue                // blue
	Magenta             // magenta
)

// FormatColor returns a color representation of given color value
func FormatColor(color int, v interface{}) string {
	return fmt.Sprintf(colorFormat, color, fmt.Sprintf("%+v", v))
}

func RedString(v interface{}) string {
	return FormatColor(Red, v)
}

func GreenString(v interface{}) string {
	return FormatColor(Green, v)
}

func YellowString(v interface{}) string {
	return FormatColor(Yellow, v)
}

func BlueString(v interface{}) string {
	return FormatColor(Blue, v)
}

func MagentaString(v interface{}) string {
	return FormatColor(Magenta, v)
}
