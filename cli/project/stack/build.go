package stack

import (
	"errors"
)

type imageBuilder func(imageName string) error

var ImageBuilders = make(map[string]imageBuilder)

func build(svcName string) error {
	if svcName == `` {
		return eachServiceDo(build)
	}
	builder, ok := ImageBuilders[svcName]
	if !ok {
		return errors.New(`no builder for ` + svcName)
	}
	imageName, err := ImageName(svcName)
	if err != nil {
		return err
	}
	return builder(imageName)
}
