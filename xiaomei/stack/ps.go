package stack

import (
	"fmt"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Ps(svcName string) error {
	var typ, name string
	if svcName != `` {
		typ, name = `service`, release.Name()+`_`+svcName
	} else {
		typ, name = `stack`, release.Name()
	}
	return cluster.Run(cmd.O{}, fmt.Sprintf(`docker %s ps --no-trunc %s`, typ, name))
}
