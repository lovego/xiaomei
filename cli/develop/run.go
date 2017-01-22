package develop

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() error {
	if err := build(); err != nil {
		return err
	}
	tail := cmd.TailFollow(
		filepath.Join(config.App.Root(), `log/app.log`),
		filepath.Join(config.App.Root(), `log/app.err`),
	)
	defer tail.Process.Kill()

	appserver.Restart(false)
	return nil
}

func build() error {
	config.Log(`building.`)
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.App.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`build failed.`)
}
