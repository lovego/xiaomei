package stack

import (
	"path"

	"github.com/bughou-go/xiaomei/config"
)

var imagesMap = make(map[string]Image)

type Image struct {
	svcName string
	image
}

type image interface {
	Prepare() error
	BuildDir() string
	Dockerfile() string
	RunPorts() []string
	RunFiles() []string
}

func RegisterImage(svcName string, img image) {
	imagesMap[svcName] = Image{svcName, img}
}

func ImageName(svcName string) string {
	return path.Join(getRegistry(), config.Name(), svcName)
}

func (i Image) Name() string {
	return ImageName(i.svcName)
}
