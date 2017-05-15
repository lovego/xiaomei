package images

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images/access"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/godoc"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/tasks"
	"github.com/lovego/xiaomei/xiaomei/images/web"
)

var imagesMap = map[string]Image{
	`app`:    Image{svcName: `app`, image: app.Image{}},
	`web`:    Image{svcName: `web`, image: web.Image{}},
	`tasks`:  Image{svcName: `tasks`, image: tasks.Image{}},
	`godoc`:  Image{svcName: `godoc`, image: godoc.Image{}},
	`access`: Image{svcName: `access`, image: access.Image{}},
	`logc`:   Image{svcName: `logc`, image: logc.Image{}, external: true},
}

func Get(svcName string) Image {
	if img, ok := imagesMap[svcName]; !ok {
		panic(`no image for: ` + svcName)
	} else {
		return img
	}
}

func Build(svcName string, pull bool) error {
	if svcName == `` {
		done := map[string]bool{}
		return eachServiceDo(func(svcName string) error {
			if imgName := conf.ImageNameOf(svcName); done[imgName] {
				return nil
			} else {
				done[imgName] = true
				return Build(svcName, pull)
			}
		})
	}
	return imagesMap[svcName].Build(pull)
}

func Push(svcName string) error {
	if svcName == `` {
		done := map[string]bool{}
		return eachServiceDo(func(svcName string) error {
			if imgName := conf.ImageNameOf(svcName); done[imgName] {
				return nil
			} else {
				done[imgName] = true
				return Push(svcName)
			}
		})
	}
	return imagesMap[svcName].Push()
}

func eachServiceDo(work func(svcName string) error) error {
	for svcName := range conf.ServiceNames() {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
