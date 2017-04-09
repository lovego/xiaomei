package app

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
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

func (i Image) EnvsForDeploy() []string {
	return []string{`GOENV=` + release.Env()}
}

func (i Image) FilesForRun() []string {
	root := release.App().Root()
	name := release.App().Name()
	return []string{
		fmt.Sprintf(`%s/%s:/home/ubuntu/%s/%s`, root, name, name, name),
		fmt.Sprintf(`%s/config:/home/ubuntu/%s/config`, root, name),
		fmt.Sprintf(`%s/views:/home/ubuntu/%s/views`, root, name),
	}
}

func (i Image) EnvsForRun() []string {
	return []string{`GODEV=true`}
}

func (i Image) CmdForRun() []string {
	return nil
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
