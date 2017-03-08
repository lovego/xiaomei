package stack

import (
	"github.com/bughou-go/xiaomei/xiaomei/images"
)

func RunImage(svcName string) error {
	if imgName, err := serviceImageName(svcName); err != nil {
		return err
	} else {
		return images.Run(svcName, imgName)
	}
}

func BuildImage(svcName string) error {
	if svcName == `` {
		return eachServiceDo(images.Build)
	}
	if imgName, err := serviceImageName(svcName); err != nil {
		return err
	} else {
		return images.Build(svcName, imgName)
	}
}

func PushImage(svcName string) error {
	if svcName == `` {
		return eachServiceDo(images.Push)
	}
	if imgName, err := serviceImageName(svcName); err != nil {
		return err
	} else {
		return images.Push(svcName, imgName)
	}
}
