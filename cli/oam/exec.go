package deploy

import (
	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"strings"
)

func Exec(args []string) {
	if len(args) == 0 {
		cli.PrintUsage()
	}
	for _, addr := range cli.MatchedServerAddrs() {
		address := config.Data.DeployUser + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `cd `+config.Data.DeployPath+`; `+strings.Join(args, ` `))
	}
}
