package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/config/servers"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"strings"
)

func Exec(serverFilter string, args []string) {
	for _, addr := range servers.MatchedAddrs(serverFilter) {
		address := config.DeployUser() + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `cd `+config.DeployPath()+`; `+strings.Join(args, ` `))
	}
}
