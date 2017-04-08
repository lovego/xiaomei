package images

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
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

func Get(svcName string) Image {
	if img, ok := imagesMap[svcName]; !ok {
		panic(`no image for: ` + svcName)
	} else {
		return img
	}
}

func Build(svcName, imgName string, pull bool) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return nil
	}
	return image.Build(imgName, pull)
}

func Push(svcName, imgName string) error {
	if _, ok := imagesMap[svcName]; !ok {
		return nil
	}
	utils.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, imgName)
	return err
}
