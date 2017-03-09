package images

import (
	"errors"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/images/app"
	"github.com/bughou-go/xiaomei/xiaomei/images/web"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/fatih/color"
)

var imagesMap = map[string]Image{
	`app`: Image{`app`, app.Image{}},
	`web`: Image{`web`, web.Image{}},
}

func Run(svcName string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Run(release.ImageNameOf(svcName))
}

func Build(svcName string) error {
	if svcName == `` {
		return eachServiceDo(Build)
	}
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Build(release.ImageNameOf(svcName))
}

func Push(svcName string) error {
	if svcName == `` {
		return eachServiceDo(Push)
	}
	utils.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, release.ImageNameOf(svcName))
	return err
}

func eachServiceDo(work func(svcName string) error) error {
	for svcName := range release.GetStack().Services {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
