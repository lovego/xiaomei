package stack

import (
	"fmt"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Ps(env, svcName string) error {
	var typ, name string
	if svcName != `` {
		typ, name = `service`, config.DeployName()+`_`+svcName
	} else {
		typ, name = `stack`, config.DeployName()
	}
	return cluster.Run(cmd.O{}, env, fmt.Sprintf(`docker %s ps %s`, typ, name))
}
