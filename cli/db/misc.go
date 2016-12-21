package db

import (
	"flag"
	"github.com/bughou-go/xiaomei/config"
	"os"
	"strings"
)

type Option struct {
	Server string
}

var ops = Option{}

func Flags() []string {
	flag.StringVar(&ops.Server, `s`, ``, `to specify server by mysql/nginx`)
	flag.Parse()
	return flag.Args()
}

func sshOptions(command string, options []string) (string, []string) {
	serverC := config.Servers.Current()
	if serverC.Addr == `` {
		addrs := getMatchedServerAddrs()
		if len(addrs) == 0 {
			os.Exit(1)
		}
		address := config.Deploy.User() + `@` + addrs[0]
		options = append([]string{`-t`, address, command}, options...)
		return `ssh`, options
	}
	return command, options
}

func getMatchedServerAddrs() []string {
	addrs := []string{}
	for _, server := range config.Servers.All() {
		if strings.Contains(server.Tasks, ops.Server) || strings.Contains(server.Addr, ops.Server) {
			if !contains(addrs, server.Addr) {
				addrs = append(addrs, server.Addr)
			}
		}
	}
	return addrs
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
