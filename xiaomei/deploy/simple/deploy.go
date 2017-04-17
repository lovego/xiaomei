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

func (d driver) Deploy(svcName string) error {
	serviceNames := []string{}
	if svcName != `` {
		serviceNames = append(serviceNames, svcName)
	} else {
		for service, toDeploy := range simpleconf.ServiceNames() {
			if toDeploy {
				serviceNames = append(serviceNames, service)
			}
		}
	}
	for _, thisSvcName := range serviceNames {
		if err := deployService(thisSvcName); err != nil {
			return err
		}
	}
	return d.Ps(svcName, true, nil)
}

func deployService(svcName string) error {
	utils.Log(color.GreenString(`deploying ` + svcName + ` service.`))

	script, err := getDeployScript(svcName)
	if err != nil {
		return err
	}

	for _, node := range cluster.Nodes() {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
