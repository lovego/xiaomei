package config

import (
	"net"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

/* basic */

func AppName() string {
	return data().AppName
}
func AppPort() string {
	return data().AppPort
}
func Env() string {
	return data().Env
}
func Domain() string {
	return data().Domain
}
func CurrentAppServer() ServerConfig {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range data().DeployServers {
		if server.AppAddr != `` {
			for _, ifcAddr := range ifcAddrs {
				if strings.HasPrefix(ifcAddr.String(), server.Addr+`/`) {
					return server
				}
			}
		}
	}
	return ServerConfig{}
}

/* for deploy */

func DeployName() string {
	return AppName() + Env()
}
func DeployRoot() string {
	return data().DeployRoot
}
func DeployPath() string {
	return filepath.Join(DeployRoot(), DeployName())
}
func DeployUser() string {
	return data().DeployUser
}
func DeployServers() []ServerConfig {
	return data().DeployServers
}

func GitAddr() string {
	return data().GitAddr
}
func GitBranch() string {
	if data().GitBranch != `` {
		return data().GitBranch
	}
	data().GitBranch, _ = cmd.Run(cmd.O{Output: true, Panic: true},
		`git`, `rev-parse`, `--abbrev-ref`, `HEAD`)
	return data().GitBranch
}

/* db config */
func Mysql() map[string]string {
	return data().Mysql
}
func Redis() map[string]string {
	return data().Redis
}
func Mongo() map[string]string {
	return data().Mongo
}
