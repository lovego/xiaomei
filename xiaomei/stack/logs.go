package stack

import (
	"fmt"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/fatih/color"
)

func Logs(svcName string, all bool) error {
	if svcName != `` {
		return serviceLog(svcName, all)
	}
	for svcName, _ = range GetStack().Services {
		if err := serviceLog(svcName, all); err != nil {
			return err
		}
	}
	return nil
}

func serviceLog(svcName string, all bool) error {
	for _, node := range cluster.GetCluster().Nodes() {
		if nodeLog(svcName, all, node) && !all {
			return nil
		}
	}
	return nil
}

func nodeLog(svcName string, all bool, node cluster.Node) bool {
	var head string
	if !all {
		head = "| head -n1"
	}
	script := fmt.Sprintf(`
ids=$(docker ps -aqf name=%s_%s.)
if test -n "$ids"; then
	echo "$ids" %s | xargs -rn1 sh -c 'echo %s $0; docker logs $0'
fi
	`,
		release.Name(), svcName, head, color.GreenString(svcName),
	)
	_, err := node.Run(cmd.O{}, script)
	return err == nil
}