package images

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
)

type Image struct {
	svcName string
	imageDriver
}

type imageDriver interface {
	PrepareForBuild() error
	BuildDir() string
	Dockerfile() string
	EnvsForDeploy() []string

	Runnable
}

type Runnable interface {
	FilesForRun() []string
	EnvsForRun() []string
	CmdForRun() []string
}

func (i Image) Build(imgName string, pull bool) error {
	if err := i.PrepareForBuild(); err != nil {
		return err
	}
	utils.Log(color.GreenString(`building ` + i.svcName + ` image.`))
	args := []string{`build`}
	if pull {
		args = append(args, `--pull`)
	}
	args = append(args, `--file=`+i.Dockerfile(), `--tag=`+imgName, `.`)
	_, err := cmd.Run(cmd.O{Dir: i.BuildDir()}, `docker`, args...)
	return err
}

func (i Image) PrepareOrBuild(imgName string) error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true},
		`docker`, `image`, `inspect`, imgName) {
		return i.PrepareForBuild()
	} else {
		return i.Build(imgName, true)
	}
}
