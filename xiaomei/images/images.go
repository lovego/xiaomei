package images

import (
	"errors"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/images/access"
	"github.com/bughou-go/xiaomei/xiaomei/images/app"
	"github.com/bughou-go/xiaomei/xiaomei/images/web"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/fatih/color"
)

var imagesMap = map[string]Image{
	`app`:    Image{`app`, app.Image{}},
	`web`:    Image{`web`, web.Image{}},
	`access`: Image{`access`, access.Image{}},
}

func Run(svcName string, ports []string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Run(ports)
}

func Build(svcName string) error {
	if svcName == `` {
		return release.EachServiceDo(Build)
	}
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Build()
}

func Push(svcName string) error {
	if svcName == `` {
		return release.EachServiceDo(Push)
	}
	utils.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, release.ImageNameOf(svcName))
	return err
}
