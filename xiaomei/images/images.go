package images

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images/access"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/web"
)

var imagesMap = map[string]Image{
	`app`:    Image{`app`, app.Image{}},
	`web`:    Image{`web`, web.Image{}},
	`access`: Image{`access`, access.Image{}},
}

func Has(svcName string) bool {
	_, ok := imagesMap[svcName]
	return ok
}

func Run(svcName string, ports []string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return nil
	}
	return image.Run(ports)
}

func Build(svcName string, pull bool) error {
	if svcName == `` {
		return eachServiceDo(func(svcName string) error {
			return Build(svcName, pull)
		})
	}
	image, ok := imagesMap[svcName]
	if !ok {
		return nil
	}
	return image.Build(pull)
}

func Push(svcName string) error {
	if svcName == `` {
		return eachServiceDo(Push)
	}
	if _, ok := imagesMap[svcName]; !ok {
		return nil
	}
	utils.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, deploy.ImageNameOf(svcName))
	return err
}

func eachServiceDo(work func(svcName string) error) error {
	for svcName := range deploy.ServiceNames() {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
