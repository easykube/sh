package sh

import (
	"testing"

	"github.com/easykube/sh/remote"
)

func TestWSByWinRm(t *testing.T) {
	conf := remote.NewConfig()
	conf.UseWinRm = true
	conf.Host = "127.0.0.1"
	conf.Port = 5985
	conf.User = "administrator"
	conf.Password = "xxxxxx"
	ws := NewWorkSpace()
	ws.Init(conf)
	ws.Run("ipconfig /all")
}
