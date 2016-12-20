package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/config/servers"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Status(serverFilter string) {
	for _, addr := range servers.MatchedAddrs(serverFilter) {
		address := config.DeployUser() + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `status `+config.DeployName()+`; ps -FC `+config.AppName())
	}
}
