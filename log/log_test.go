package log

import (
	"errors"
	"testing"
)

func TestConsoleLog(t *testing.T) {
	Info("info")
	Error(errors.New("error"))
	Warn(struct {
		Name string
	}{
		Name: "warn",
	})
	Debug(0x666)
}
