package images

import (
	"github.com/lovego/xiaomei/services/images/app"
	"github.com/lovego/xiaomei/services/images/logc"
	"github.com/lovego/xiaomei/services/images/web"
)

type Image struct {
	svcName string
	image   interface{}
}

func Get(svcName string) Image {
	switch svcName {
	case "app":
		return Image{svcName: `app`, image: app.Image{}}
	case "web":
		return Image{svcName: `web`, image: web.Image{}}
	case "logc":
		return Image{svcName: `logc`, image: logc.Image{}}
	default:
		panic(`no image for: ` + svcName)
	}
}

// the 4 interfaces, image driver choose to implement.

// 1. port env variable name
func (i Image) PortEnvVar() string {
	if img, ok := i.image.(interface {
		PortEnvVar() string
	}); ok {
		return img.PortEnvVar()
	}
	return ``
}

// 2. default port number
func (i Image) DefaultPort() uint16 {
	if img, ok := i.image.(interface {
		DefaultPort() uint16
	}); ok {
		return img.DefaultPort()
	}
	return 0
}

// 3. flags for run
func (i Image) FlagsForRun() []string {
	if img, ok := i.image.(interface {
		OptionsForRun() []string
	}); ok {
		return img.OptionsForRun()
	}
	return nil
}

// 4. prepare files for build
func (i Image) prepare(env string, flags []string) error {
	if img, ok := i.image.(interface {
		Prepare(env string, flags []string) error
	}); ok {
		return img.Prepare(env, flags)
	}
	return nil
}
