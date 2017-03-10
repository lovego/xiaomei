package images

import (
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/fatih/color"
)

type Image struct {
	svcName string
	imageDriver
}

type imageDriver interface {
	PrepareForBuild() error
	BuildDir() string
	Dockerfile() string
	RunFiles() []string
	RunCmd() string
}

func (i Image) Build() error {
	if err := i.PrepareForBuild(); err != nil {
		return err
	}
	utils.Log(color.GreenString(`building ` + i.svcName + ` image.`))
	_, err := cmd.Run(cmd.O{Dir: i.BuildDir()}, `docker`, `build`,
		`--file=`+i.Dockerfile(), `--tag=`+release.ImageNameOf(i.svcName), `.`,
	)
	return err
}
