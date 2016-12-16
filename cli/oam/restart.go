package oam

import (
	"fmt"
	"os"
	"strings"

	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Restart() {
	addrs := cli.MatchedServerAddrs()
	for _, addr := range addrs {
		restartAppServer(config.DeployUser() + `@` + addr)
	}
	fmt.Printf("restart %d server!\n", len(addrs))
}

func restartAppServer(address string) {
	color.Cyan(address)

	command := fmt.Sprintf(
		`sudo stop %s; sudo start %s`,
		config.DeployName(), config.DeployName(),
	)
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `ssh`, `-t`, address, command)

	if strings.Contains(output, `start/running,`) {
		fmt.Printf("restart %s ok.\n", config.DeployName())
	} else {
		fmt.Printf("restart %s failed.\n", config.DeployName())
		os.Exit(1)
	}
}
