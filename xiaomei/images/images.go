package images

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/godoc"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/tasks"
	"github.com/lovego/xiaomei/xiaomei/images/web"
)

var imagesMap = map[string]Image{
	`app`:   Image{svcName: `app`, image: app.Image{}},
	`tasks`: Image{svcName: `tasks`, image: tasks.Image{}},
	`web`:   Image{svcName: `web`, image: web.Image{}},
	`logc`:  Image{svcName: `logc`, image: logc.Image{}},
	`godoc`: Image{svcName: `godoc`, image: godoc.Image{}},
}

func Get(svcName string) Image {
	if img, ok := imagesMap[svcName]; !ok {
		panic(`no image for: ` + svcName)
	} else {
		return img
	}
}

func Build(env, svcName string, pull bool) error {
	if svcName == `` {
		done := map[string]bool{}
		return eachServiceDo(func(svcName string) error {
			if imgName := conf.ImageNameOf(svcName); done[imgName] {
				return nil
			} else {
				done[imgName] = true
				return Build(env, svcName, pull)
			}
		})
	}
	return imagesMap[svcName].Build(env, pull)
}

func Push(env, svcName string) error {
	if svcName == `` {
		done := map[string]bool{}
		return eachServiceDo(func(svcName string) error {
			if imgName := conf.ImageNameOf(svcName); done[imgName] {
				return nil
			} else {
				done[imgName] = true
				return Push(env, svcName)
			}
		})
	}
	return imagesMap[svcName].Push(env)
}

func eachServiceDo(work func(svcName string) error) error {
	for _, svcName := range conf.ServiceNames() {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
