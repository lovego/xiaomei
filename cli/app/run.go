package app

import (
	"errors"

	"github.com/bughou-go/xiaomei/cli/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() error {
	if err := buildBinary(); err != nil {
		return err
	}

	appImage, err := stack.ServiceImage(`app`)
	if err != nil {
		return err
	}
	_, err = cmd.Run(cmd.O{}, `docker`,
		`run`, `--name=`+config.DeployName(), `-it`, `--rm`, `--network=host`,
		`-v`, config.Root()+`:/home/ubuntu/appserver`,
		appImage,
	)
	return err
}

func buildBinary() error {
	config.Log(`building.`)
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`build failed.`)
}
