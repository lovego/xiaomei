package images

import (
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/images/app"
	"github.com/bughou-go/xiaomei/xiaomei/images/web"
	"github.com/fatih/color"
)

var imagesMap = map[string]Image{
	`app`: Image{`app`, app.Image{}},
	`web`: Image{`web`, web.Image{}},
}

func Run(svcName, imgName string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Run(imgName)
}

func Build(svcName, imgName string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Build(imgName)
}

func Push(svcName, imgName string) error {
	config.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, imgName)
	return err
}
