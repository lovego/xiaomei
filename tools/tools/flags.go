package tools

import (
	"flag"
	"github.com/bughou-go/xiaomei/config"
	"strings"
)

type Options struct {
	Server string
}

var options = Options{}

func Flags() []string {
	flag.StringVar(&options.Server, `s`, ``, `to specify server by mysql/nginx`)
	flag.Parse()
	return flag.Args()
}

func MatchedServers() []config.ServerConfig {
	matched := []config.ServerConfig{}
	for _, server := range config.Data.DeployServers {
		if strings.Contains(server.Tasks, options.Server) ||
			strings.Contains(server.Addr, options.Server) {
			matched = append(matched, server)
		}
	}
	return matched
}

func MatchedServerAddrs() []string {
	addrs := []string{}
	for _, server := range config.Data.DeployServers {
		if strings.Contains(server.Tasks, options.Server) ||
			strings.Contains(server.Addr, options.Server) {
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
