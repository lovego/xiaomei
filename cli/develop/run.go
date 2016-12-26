package develop

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() {
	if err := build(); err != nil {
		return
	}
	config.Log(`starting.`)
	app, err := cmd.Start(cmd.O{}, filepath.Join(config.App.Root(), config.App.Name()))
	if err != nil {
		return
	}
	appserver.WaitPort(os.Getpid(), app.Process.Pid)
	app.Wait()
}

func build() error {
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.App.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`build failed.`)
}
