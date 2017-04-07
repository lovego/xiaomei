package host

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
)

func (d driver) Deploy(svcName string) error {
	serviceNames := []string{}
	if svcName != `` {
		serviceNames = append(serviceNames, svcName)
	} else {
		services := getRelease()
		if _, ok := services[`app`]; ok {
			serviceNames = append(serviceNames, `app`)
		}
		if _, ok := services[`web`]; ok {
			serviceNames = append(serviceNames, `web`)
		}
	}
	for _, svcName := range serviceNames {
		if err := deployService(svcName); err != nil {
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
