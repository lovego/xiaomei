package images

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy"
)

type Image struct {
	svcName string
	imageDriver
}

type imageDriver interface {
	PrepareForBuild() error
	BuildDir() string
	Dockerfile() string
	deploy.Runnable
}

func (i Image) Build(pull bool) error {
	if err := i.PrepareForBuild(); err != nil {
		return err
	}
	utils.Log(color.GreenString(`building ` + i.svcName + ` image.`))
	args := []string{`build`}
	if pull {
		args = append(args, `--pull`)
	}
	args = append(args, `--file=`+i.Dockerfile(), `--tag=`+deploy.ImageNameOf(i.svcName), `.`)
	_, err := cmd.Run(cmd.O{Dir: i.BuildDir()}, `docker`, args...)
	return err
}

func (i Image) PrepareOrBuild(i Image) error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true},
		`docker`, `image`, `inspect`, deploy.ImageNameOf(i.svcName)) {
		return i.PrepareForBuild()
	} else {
		return i.Build(true)
	}
}
