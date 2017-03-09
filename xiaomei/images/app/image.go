package app

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
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
	return release.App().Root()
}

func (i Image) Dockerfile() string {
	return `Dockerfile`
}

func (i Image) RunPorts() []string {
	return []string{`3000:3000`}
}

func (i Image) RunFiles() []string {
	return []string{
		release.App().Root() + `:/home/ubuntu/` + release.App().Name(),
	}
}

func buildBinary() error {
	utils.Log(color.GreenString(`building app binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.App().Root(), "../.."),
		Env: []string{`GOBIN=` + release.App().Root()},
	}, `go`, `install`) {
		return nil
	}
	return errors.New(`building app binary failed.`)
}
