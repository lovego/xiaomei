package deploy

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"strings"
)

func Restart() {
	addrs := tools.MatchedServerAddrs()
	for _, addr := range addrs {
		restartAppServer(config.Data.DeployUser + `@` + addr)
	}
	fmt.Printf("restart %d server!\n", len(addrs))
}

func restartAppServer(address string) {
	color.Cyan(address)

	command := fmt.Sprintf(
		`sudo stop %s; sudo start %s`,
		config.Data.DeployName, config.Data.DeployName,
	)
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `ssh`, address, command)

	if strings.Contains(output, `start/running,`) {
		fmt.Printf("restart %s ok.\n", config.Data.DeployName)
	} else {
		fmt.Printf("restart %s failed.\n", config.Data.DeployName)
		os.Exit(1)
	}
}
