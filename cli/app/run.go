package app

import (
	"errors"

	"github.com/bughou-go/xiaomei/cli/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() error {
	if err := BuildBinary(); err != nil {
		return err
	}

	var image string
	if svc, err := stack.GetService(`app`); err != nil {
		return err
	} else if image, err = svc.GetImage(); err != nil {
		return err
	}

	_, err := cmd.Run(cmd.O{}, `docker`,
		`run`, `--name=`+config.DeployName(), `-it`, `--rm`, `--network=host`,
		`-v`, config.Root()+`:/home/ubuntu/appserver`,
		image,
	)
	return err
}

func Build() error {
	if err := BuildBinary(); err != nil {
		return err
	}
	if err := Spec(``); err != nil {
		return err
	}
	// Assets(nil)

	return BuildImage()
}

func BuildBinary() error {
	config.Log(`building binary.`)
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`building binary failed.`)
}

func BuildImage() error {
	if svc, err := stack.GetService(`app`); err != nil {
		return err
	} else if err = svc.BuildImage(); err != nil {
		return err
	}
	return nil
}
