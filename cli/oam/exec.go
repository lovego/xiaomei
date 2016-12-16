package oam

import (
	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"strings"
)

func Exec(args []string) {
	for _, addr := range cli.MatchedServerAddrs() {
		address := config.DeployUser() + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `cd `+config.DeployPath()+`; `+strings.Join(args, ` `))
	}
}
