package oam

import (
	"fmt"
	"os"
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Restart(serverFilter string) {
	addrs := config.Servers.MatchedAddrs(serverFilter)
	for _, addr := range addrs {
		restartAppServer(config.Deploy.User() + `@` + addr)
	}
	fmt.Printf("restart %d server!\n", len(addrs))
}

func restartAppServer(address string) {
	color.Cyan(address)

	command := fmt.Sprintf(
		`sudo stop %s; sudo start %s`,
		config.Deploy.Name(), config.Deploy.Name(),
	)
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `ssh`, `-t`, address, command)

	if strings.Contains(output, `start/running,`) {
		fmt.Printf("restart %s ok.\n", config.Deploy.Name())
	} else {
		fmt.Printf("restart %s failed.\n", config.Deploy.Name())
		os.Exit(1)
	}
}
