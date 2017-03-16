package stack

import (
	"fmt"
	"strings"

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
		if has, err := nodeLog(svcName, all, node); err != nil {
			return err
		} else if !all && has {
			return nil
		}
	}
	return nil
}

func nodeLog(svcName string, all bool, node cluster.Node) (bool, error) {
	script := fmt.Sprintf(`docker ps -aqf name=%s_%s.`, release.Name(), svcName)
	if output, err := node.Run(cmd.O{Output: true}, script); err != nil {
		return false, err
	} else {
		containerIds := strings.Split(output, "\n")
		return len(containerIds) > 0, containersLog(svcName, all, node, containerIds)
	}
}

func containersLog(svcName string, all bool, node cluster.Node, containerIds []string) error {
	for _, containerId := range containerIds {
		println(color.GreenString(`%s %s:`, svcName, containerId))

		if _, err := node.Run(cmd.O{}, fmt.Sprintf(`docker logs %s`, containerId)); err != nil {
			return err
		}
		if !all {
			break
		}
	}
	return nil
}
