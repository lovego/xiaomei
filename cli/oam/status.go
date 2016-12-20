package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Status(serverFilter string) {
	for _, addr := range config.Servers.MatchedAddrs(serverFilter) {
		address := config.Deploy.User() + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `status `+config.Deploy.Name()+`; ps -FC `+config.App.Name())
	}
}
