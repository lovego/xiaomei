package stack

import (
	"fmt"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Push(svcName string) error {
	if svcName == `` {
		return eachServiceDo(Push)
	}
	config.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, ImageName(svcName))
	return err
}

func Ps(svcName string) error {
	var typ, name string
	if svcName != `` {
		typ, name = `service`, config.Name()+`_`+svcName
	} else {
		typ, name = `stack`, config.Name()
	}
	return cluster.Run(cmd.O{}, config.Env(),
		fmt.Sprintf(`docker %s ps --no-trunc %s`, typ, name))
}
