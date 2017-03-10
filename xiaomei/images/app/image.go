package app

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/fatih/color"
)

type Image struct {
}

func (i Image) PrepareForBuild() error {
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

func (i Image) RunFiles() []string {
	root := release.App().Root()
	name := release.App().Name()
	return []string{
		fmt.Sprintf(`%s/%s:/home/ubuntu/%s/%s`, root, name, name, name),
		fmt.Sprintf(`%s/config:/home/ubuntu/%s/config`, root, name),
		fmt.Sprintf(`%s/views:/home/ubuntu/%s/views`, root, name),
	}
}

func (i Image) RunCmd() string {
	return ``
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
