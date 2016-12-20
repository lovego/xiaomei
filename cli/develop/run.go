package develop

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() error {
	if err := build(); err != nil {
		return err
	}
	if cmd.Ok(cmd.O{}, filepath.Join(config.App.Root(), config.AppName())) {
		return nil
	}
	return errors.New(`run failed.`)
}

func build() error {
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.App.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`build failed.`)
}
