package regia

import (
	"errors"
	"testing"
)

func TestConsoleLog(t *testing.T) {
	Logger.Info("info")
	Logger.Error(errors.New("error"))
	Logger.Warn(struct {
		Name string
	}{
		Name: "warn",
	})
	Logger.Debug(0x666)
}
