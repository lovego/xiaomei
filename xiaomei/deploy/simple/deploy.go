package simple

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
)

var Driver driver

type driver struct{}

func (d driver) Deploy(svcName, feature string) error {
	svcs := getServices(svcName)
	psScript := getPsScript(svcName, true)
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
		return simpleconf.ServiceNames()
	} else {
		return []string{svcName}
	}
}

func getNodeServices(svcNames []string, node cluster.Node) []string {
	svcs := []string{}
	for _, svcName := range svcNames {
		service := simpleconf.GetService(svcName)
		if node.Match(service.Nodes) {
			svcs = append(svcs, svcName)
		}
	}
	return svcs
}
