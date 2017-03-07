package app

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func init() {
	stack.RegisterImage(`app`, appImage{})
}

type appImage struct {
}

func (app appImage) Prepare() error {
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

func (app appImage) BuildDir() string {
	return config.Root()
}

func (app appImage) Dockerfile() string {
	return `Dockerfile`
}

func (app appImage) RunMapping() []string {
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
