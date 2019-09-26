package images

import (
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images/app"
	"github.com/lovego/xiaomei/services/images/logc"
	"github.com/lovego/xiaomei/services/images/web"
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

func Build(svcName, env, tag string, pull bool) error {
	return imagesDo(svcName, env, func(img Image) error {
		return img.build(env, tag, pull)
	})
}

func Push(svcName, env, tag string) error {
	return imagesDo(svcName, env, func(img Image) error {
		return img.push(env, tag)
	})
}

func List(svcName, env string) error {
	return imagesDo(svcName, env, func(img Image) error {
		return img.list(env)
	})
}

func imagesDo(svcName, env string, work func(Image) error) error {
	if svcName == `` {
		for _, svcName := range release.ServiceNames(env) {
			if err := work(Get(svcName)); err != nil {
				return err
			}
		}
		return nil
	} else {
		return work(Get(svcName))
	}
}
