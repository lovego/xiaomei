package deploy

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func deploy(svcName, feature string) error {
	svcs := getServices(svcName)
	psScript := fmt.Sprintf(`watch docker ps -f name=%s_%s`, release.DeployName(), svcName)
	for _, node := range cluster.Nodes(feature) {
		if err := deployNode(svcs, node, psScript); err != nil {
			return err
		}
	}
	return nil
}

func deployNode(svcs []string, node cluster.Node, psScript string) error {
	nodeSvcs := getNodeServices(svcs, node)
	if len(nodeSvcs) == 0 {
		return nil
	}
	utils.Log(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(nodeSvcs)
	if err != nil {
		return err
	}
	_, err = node.Run(cmd.O{}, deployScript+psScript)
	return err
}

func getServices(svcName string) []string {
	if svcName == `` {
		return conf.ServiceNames()
	} else {
		return []string{svcName}
	}
}

func getNodeServices(svcNames []string, node cluster.Node) []string {
	svcs := []string{}
	for _, svcName := range svcNames {
		service := conf.GetService(svcName)
		if node.Match(service.Nodes) {
			svcs = append(svcs, svcName)
		}
	}
	return svcs
}
