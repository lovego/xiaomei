package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"strings"
)

func Exec(serverFilter string, args []string) {
	for _, addr := range config.Servers.MatchedAddrs(serverFilter) {
		address := config.Deploy.User() + `@` + addr
		color.Cyan(address)
		cmd.Run(cmd.O{}, `ssh`, `-t`, address, `cd `+config.Deploy.Path()+`; `+strings.Join(args, ` `))
	}
}
