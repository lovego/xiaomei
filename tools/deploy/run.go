package deploy

import (
	"github.com/fatih/color"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"strings"
)

func Run(args []string) {
	if len(args) == 0 {
		tools.PrintUsage()
	}
	for _, addr := range tools.GetMatchedServerAddrs() {
		address := config.Data.DeployUser + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, address, `cd `+config.Data.DeployPath+`; `+strings.Join(args, ` `))
	}
}
