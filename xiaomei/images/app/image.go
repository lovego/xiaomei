package app

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

type Image struct {
}

func (i Image) Prepare() error {
	if err := buildBinary(); err != nil {
		return err
	}
	/*
		if err :=	Assets(nil); err != nil {
			return err
		}
	*/
	return nil
}

func (i Image) BuildDir() string {
	return config.Root()
}

func (i Image) Dockerfile() string {
	return `Dockerfile`
}

func (i Image) RunPorts() []string {
	return []string{`3000:3000`}
}

func (i Image) RunFiles() []string {
	return []string{
		config.Root() + `:/home/ubuntu/` + config.Name(),
	}
}

func buildBinary() error {
	config.Log(color.GreenString(`building app binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(config.Root(), "../.."),
		Env: []string{`GOBIN=` + config.Root()},
	}, `go`, `install`) {
		return nil
	}
	return errors.New(`building app binary failed.`)
}
