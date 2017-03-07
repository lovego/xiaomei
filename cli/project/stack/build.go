package stack

import (
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Build(svcName string) error {
	if svcName == `` {
		return eachServiceDo(Build)
	}
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Build()
}

func (i Image) Build() error {
	if err := i.Prepare(); err != nil {
		return err
	}
	config.Log(color.GreenString(`building ` + i.svcName + ` image.`))
	_, err := cmd.Run(cmd.O{Dir: i.BuildDir()}, `docker`, `build`,
		`--file=`+i.Dockerfile(), `--tag=`+i.Name(), `.`,
	)
	return err
}
