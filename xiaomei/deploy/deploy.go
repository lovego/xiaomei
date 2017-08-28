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

func deploy(env, svcName, feature string) error {
	svcs := getServices(env, svcName)
	psScript := fmt.Sprintf(`watch docker ps -f name=%s`, release.ServiceName(env, svcName))
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if err := deployNode(env, svcs, node, psScript); err != nil {
			return err
		}
	}
	return nil
}

func deployNode(env string, svcs []string, node cluster.Node, psScript string) error {
	nodeSvcs := getNodeServices(env, svcs, node)
	if len(nodeSvcs) == 0 {
		return nil
	}
	utils.Log(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(env, nodeSvcs)
	if err != nil {
		return err
	}
	_, err = node.Run(cmd.O{}, deployScript+psScript)
	return err
}

func getServices(env, svcName string) []string {
	if svcName == `` {
		return conf.ServiceNames(env)
	} else {
		return []string{svcName}
	}
}

func getNodeServices(env string, svcNames []string, node cluster.Node) []string {
	svcs := []string{}
	for _, svcName := range svcNames {
		service := conf.GetService(env, svcName)
		if node.Match(service.Nodes) {
			svcs = append(svcs, svcName)
		}
	}
	return svcs
}
