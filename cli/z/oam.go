package oam

import (
	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"strings"
)

func Status(serverFilter string, args []string) error {
	const format = `table {{.ID }}\t{{.Names}}\t{{.Status}}\t{{.CreatedAt}}\t{{.Image}}\t{{.Command}}`
	return run(serverFilter,
		`docker ps -af name=`+config.Deploy.Name()+` --format='`+format+`'
		pid=$(docker inspect --type=container -f '{{ .State.Pid }}' `+config.Deploy.Name()+`)
		test -n "$pid" && echo && ps -F $pid; true`)
}

func Restart(serverFilter string, args []string) error {
	if config.IsLocalEnv() {
		appserver.Restart(true)
		return nil
	} else {
		return run(serverFilter, `cd `+config.Deploy.Path()+`; xiaomei restart`)
	}
}

func Shell(serverFilter string, args []string) error {
	return runOnce(serverFilter, `cd `+config.Deploy.Path()+`; bash`)
}

func Exec(serverFilter string, args []string) error {
	return run(serverFilter, `cd `+config.Deploy.Path()+`; `+strings.Join(args, ` `))
}

func run(filter, shellScript string) error {
	if config.IsLocalEnv() {
		_, err := cmd.Run(cmd.O{}, `sh`, `-c`, shellScript)
		return err
	}
	for _, server := range config.Servers.Matched2(filter, `appserver`) {
		addr := server.SshAddr()
		color.Cyan(addr)
		_, err := cmd.Run(cmd.O{}, `ssh`, `-t`, addr, shellScript)
		if err != nil {
			return err
		}
	}
	return nil
}

func runOnce(filter, shellScript string) error {
	if config.IsLocalEnv() {
		_, err := cmd.Run(cmd.O{}, `sh`, `-c`, shellScript)
		return err
	}
	for _, server := range config.Servers.Matched2(filter, `appserver`) {
		addr := server.SshAddr()
		color.Cyan(addr)
		_, err := cmd.Run(cmd.O{}, `ssh`, `-t`, addr, shellScript)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
