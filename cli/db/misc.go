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
		options = append([]string{`-t`, servers[0].SshAddr(), command}, options...)
		return `ssh`, options
	}
	return command, options
}
