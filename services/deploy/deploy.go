package deploy

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
	"github.com/lovego/xiaomei/services/oam"
)

type Deploy struct {
	svcName string
	images.Build
	images.Push
	alwaysPush bool

	filter  string
	noWatch bool
}

func (d Deploy) start() error {
	if d.Build.Tag == `` {
		d.Build.Tag = release.TimeTag(d.Build.Env)
		d.Push.Tag = d.Build.Tag
		if err := d.createImages(); err != nil {
			return err
		}
	}
	if err := d.run(); err != nil {
		return err
	}
	return nil
}

func (d Deploy) createImages() error {
	if err := d.Build.Run(d.svcName); err != nil {
		return err
	}
	if !d.alwaysPush {
		if ok, err := release.GetCluster(d.Build.Env).IsLocalHost(); err != nil {
			return err
		} else if ok {
			return nil
		}
	}

	return d.Push.Run(d.svcName)
}

func (d Deploy) run() error {
	psScript := fmt.Sprintf(` docker ps -f name=^/%s`, release.ServiceName(d.svcName, d.Build.Env))
	if !d.noWatch {
		psScript = oam.WatchCmd() + psScript
	}
	for _, node := range release.GetCluster(d.Build.Env).GetNodes(d.filter) {
		if svcs := node.Services(d.Build.Env, d.svcName); len(svcs) > 0 {
			if err := d.runNode(svcs, node, psScript); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d Deploy) runNode(svcs []string, node release.Node, psScript string) error {
	log.Println(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(svcs, d.Build.Env, d.Build.Tag)
	if err != nil {
		return err
	}
	if ok, err := node.IsLocalHost(); err != nil {
		return err
	} else if !ok {
		deployScript = d.Push.DockerLogin.BashCommand(d.Build.Env, svcs) + "\n" + deployScript
	}
	_, err = node.Run(cmd.O{}, deployScript+"\n"+psScript)
	return err
}
