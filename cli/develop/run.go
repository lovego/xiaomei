package develop

import (
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() {
	if !build() {
		return
	}
	cmd.Run(cmd.O{}, filepath.Join(config.Root(), config.Data().AppName))
}

func build() bool {
	env := append(os.Environ(), `GOBIN=`+config.Root())
	return cmd.Ok(cmd.O{Env: env}, `go`, `install`)
}
