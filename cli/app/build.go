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
	stack.ImageBuilders[`app`] = buildImage
}

func buildImage(imageName string) error {
	if err := buildBinary(); err != nil {
		return err
	}
	/*
		if err :=	Assets(nil); err != nil {
			return err
		}
	*/
	config.Log(color.GreenString(`building app image.`))
	dir := filepath.Join(config.Root(), `../img-app`)
	_, err := cmd.Run(cmd.O{Dir: dir}, `docker`, `build`, `--tag=`+imageName, `.`)
	return err
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
