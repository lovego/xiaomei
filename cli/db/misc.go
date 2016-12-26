package db

import (
	"os"

	"github.com/bughou-go/xiaomei/config"
)

func sshOptions(command string, options []string) (string, []string) {
	server := config.Servers.CurrentAppServer()
	if server == nil {
		addrs := config.Servers.MatchedAddrs(``)
		if len(addrs) == 0 {
			os.Exit(1)
		}
		address := config.Deploy.User() + `@` + addrs[0]
		options = append([]string{`-t`, address, command}, options...)
		return `ssh`, options
	}
	return command, options
}
