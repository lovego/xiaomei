package deploy

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/oam"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func deploy(svcName, env, timeTag, feature string) error {
	svcs := getServices(env, svcName)
	ha := expectHighAvailable(env, svcs)
	psScript := fmt.Sprintf(oam.WatchCmd()+` docker ps -f name=%s`, release.ServiceName(svcName, env))
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if ha {
			if err := access.SetupNginx(env, "", "", node.Addr); err != nil {
				return err
			}
			time.Sleep(time.Second) // wait for nginx reloading finished.
		}
		if err := deployNode(svcs, env, timeTag, node, psScript); err != nil {
			return err
		}
	}
	if ha {
		return access.SetupNginx(env, "", "", "")
	}
	return nil
}

func expectHighAvailable(env string, svcs []string) bool {
	if len(cluster.Get(env).GetNodes("")) < 2 {
		return false
	}
	for _, svcName := range svcs {
		if svcName == "app" || svcName == "web" {
			return true
		}
	}
	return false
}

func deployNode(svcs []string, env, timeTag string, node cluster.Node, psScript string) error {
	nodeSvcs := getNodeServices(svcs, env, node)
	if len(nodeSvcs) == 0 {
		return nil
	}
	log.Println(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(nodeSvcs, env, timeTag)
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

func getNodeServices(svcNames []string, env string, node cluster.Node) []string {
	svcs := []string{}
	for _, svcName := range svcNames {
		service := conf.GetService(svcName, env)
		if node.Match(service.Nodes) {
			svcs = append(svcs, svcName)
		}
	}
	return svcs
}
