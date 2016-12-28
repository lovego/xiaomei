package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Shell(serverFilter string) {
	for _, server := range config.Servers.MatchedAppserver(serverFilter) {
		cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, server.SshAddr(), `cd `+config.Deploy.Path()+`; bash`)
		return
	}
}
