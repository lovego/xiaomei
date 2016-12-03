package setup

import (
	"fmt"
	"path"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func SetupHosts() {
	cmd.Run(cmd.O{Panic: true}, path.Join(config.Root(), `config/shell/setup-hosts.sh`))

	fmt.Println(`setup hosts ok.`)
}
