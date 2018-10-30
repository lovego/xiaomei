package images

import (
	"github.com/lovego/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/images/app"
	"github.com/lovego/xiaomei/images/logc"
	"github.com/lovego/xiaomei/images/web"
)

var imagesMap = map[string]Image{
	`app`:  {svcName: `app`, image: app.Image{}},
	`web`:  {svcName: `web`, image: web.Image{}},
	`logc`: {svcName: `logc`, image: logc.Image{}},
}

func Get(svcName string) Image {
	if img, ok := imagesMap[svcName]; !ok {
		panic(`no image for: ` + svcName)
	} else {
		return img
	}
}

func Build(svcName, env, timeTag string, pull bool) error {
	return imagesDo(svcName, env, func(img Image) error {
		return img.build(env, timeTag, pull)
	})
}

func Push(svcName, env, timeTag string) error {
	return imagesDo(svcName, env, func(img Image) error {
		return img.push(env, timeTag)
	})
}

func imagesDo(svcName, env string, work func(Image) error) error {
	if svcName == `` {
		for _, svcName := range conf.ServiceNames(env) {
			if err := work(Get(svcName)); err != nil {
				return err
			}
		}
		return nil
	} else {
		return work(Get(svcName))
	}
}
