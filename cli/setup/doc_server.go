package setup

import (
	"path/filepath"
	"syscall"

	"github.com/bughou-go/xiaomei/cli/install"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func DocServer() error {
	cmd.Run(cmd.O{}, `killall`, `-q`, `godoc`)

	if cmd.Fail(cmd.O{NoStdout: true}, `which`, `godoc`) {
		if err := install.AptGet(`golang-go.tools`); err != nil {
			return err
		}
	}

	deployRoot := config.Data().DeployRoot
	if _, err := cmd.Run(cmd.O{}, `ln`, `-sf`, `.`, filepath.Join(deployRoot, `src`)); err != nil {
		return err
	}

	_, err := cmd.Start(cmd.O{
		Env:  []string{`GOPATH=` + deployRoot},
		Attr: &syscall.SysProcAttr{Setpgid: true},
	}, `godoc`, `-http=:1234`)
	return err
}
