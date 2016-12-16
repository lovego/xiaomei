package oam

import (
	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Status() {
	for _, addr := range cli.MatchedServerAddrs() {
		address := config.DeployUser() + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `status `+config.DeployName()+`; ps -FC appserver`)
	}
}
