package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Shell(serverFilter string) {
	for _, addr := range config.Servers.MatchedAddrs(serverFilter) {
		address := config.Deploy.User() + `@` + addr
		cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, address, `cd `+config.Deploy.Path()+`; bash`)
		return
	}
}
