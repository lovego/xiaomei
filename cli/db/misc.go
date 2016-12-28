package db

import (
	"os"

	"github.com/bughou-go/xiaomei/config"
)

func sshOptions(command string, options []string) (string, []string) {
	server := config.Servers.CurrentAppServer()
	if server == nil {
		servers := config.Servers.MatchedAppserver(``)
		if len(servers) == 0 {
			os.Exit(1)
		}
		address := config.Deploy.User() + `@` + servers[0].SshAddr()
		options = append([]string{`-t`, address, command}, options...)
		return `ssh`, options
	}
	return command, options
}
