package stack

import (
	"fmt"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Ps(svcName string, options []string) error {
	var script string
	if svcName != `` {
		script = fmt.Sprintf(`docker service ps %s_%s`, release.Name(), svcName)
	} else {
		script = fmt.Sprintf(`docker stack ps %s`, release.Name())
	}
	if len(options) > 0 {
		script += ` ` + strings.Join(options, ` `)
	}

	return cluster.Run(cmd.O{}, script)
}
