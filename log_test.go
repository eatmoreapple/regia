package regia

import (
	"errors"
	"testing"
)

func TestConsoleLog(t *testing.T) {
	Log.Info("info")
	Log.Error(errors.New("error"))
	Log.Warn(struct {
		Name string
	}{
		Name: "warn",
	})
	Log.Debug(0x666)
}
