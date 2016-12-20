package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/config/servers"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Shell(serverFilter string) {
	for _, addr := range servers.MatchedAddrs(serverFilter) {
		address := config.DeployUser() + `@` + addr
		cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, address, `cd `+config.DeployPath()+`; bash`)
		return
	}
}
