package deploy

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/access"
	"github.com/lovego/xiaomei/cluster"
	"github.com/lovego/xiaomei/oam"
	"github.com/lovego/xiaomei/release"
)

func deploy(svcName, env, timeTag, feature string) error {
	psScript := fmt.Sprintf(oam.WatchCmd()+` docker ps -f name=^/%s`, release.ServiceName(svcName, env))
	var hasAccess bool
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if svcs := node.Services(env, svcName); len(svcs) > 0 {
			if access.HasAccess(svcs) {
				if err := access.SetupNginx(env, "", "", node.Addr); err != nil {
					return err
				}
				hasAccess = true
				time.Sleep(time.Second) // wait for nginx reloading finished.
			}
			if err := deployNode(svcs, env, timeTag, node, psScript); err != nil {
				return err
			}
		}
	}
	if hasAccess {
		return access.SetupNginx(env, "", "", "")
	}
	return nil
}

func deployNode(svcs []string, env, timeTag string, node cluster.Node, psScript string) error {
	log.Println(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(svcs, env, timeTag)
	if err != nil {
		return err
	}
	_, err = node.Run(cmd.O{}, deployScript+psScript)
	return err
}
