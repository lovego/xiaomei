package deploy

import (
	"github.com/fatih/color"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Status() {
	for _, addr := range tools.MatchedServerAddrs() {
		address := config.Data.DeployUser + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `status `+config.Data.DeployName+`; ps -FC appserver`)
	}
}
