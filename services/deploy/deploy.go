package deploy

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/access"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
	"github.com/lovego/xiaomei/services/oam"
)

type Deploy struct {
	svcName string
	images.Build
	noPushImageIfLocal bool

	filter                string
	beforeScript          string
	noBeforeScriptOnLocal bool
	noWatch               bool
}

func (d Deploy) start() error {
	if d.beforeScript != `` && !d.noBeforeScriptOnLocal {
		if _, err := cmd.Run(cmd.O{}, `bash`, `-c`, d.beforeScript); err != nil {
			return err
		}
	}
	if d.Tag == `` {
		d.Tag = release.TimeTag(d.Env)
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
	if d.noPushImageIfLocal {
		if ok, err := release.GetCluster(d.Env).IsLocalHost(); err != nil {
			return err
		} else if ok {
			return nil
		}
	}
	return images.Push(d.svcName, d.Env, d.Tag)
}

func (d Deploy) run() error {
	psScript := fmt.Sprintf(` docker ps -f name=^/%s`, release.ServiceName(d.svcName, d.Env))
	if !d.noWatch {
		psScript = oam.WatchCmd() + psScript
	}
	expectHighAvailable := len(release.GetCluster(d.Env).GetNodes("")) >= 2
	var recoverAccess bool
	for _, node := range release.GetCluster(d.Env).GetNodes(d.filter) {
		if svcs := node.Services(d.Env, d.svcName); len(svcs) > 0 {
			if access.HasAccess(svcs) {
				if expectHighAvailable {
					if err := access.SetupNginx(d.Env, "", node.Addr); err != nil {
						return err
					}
					time.Sleep(time.Second) // wait for nginx reloading finished.
				}
				recoverAccess = true
			}
			if err := d.runNode(svcs, node, psScript); err != nil {
				return err
			}
		}
	}
	if recoverAccess {
		return access.SetupNginx(d.Env, "", "")
	}
	return nil
}

func (d Deploy) runNode(svcs []string, node release.Node, psScript string) error {
	log.Println(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(svcs, d.Env, d.Tag)
	if err != nil {
		return err
	}
	if d.beforeScript != `` {
		if d.noBeforeScriptOnLocal {
			deployScript = d.beforeScript + "\n" + deployScript
		} else if ok, err := node.IsLocalHost(); err != nil {
			return err
		} else if !ok {
			deployScript = d.beforeScript + "\n" + deployScript
		}
	}

	_, err = node.Run(cmd.O{}, deployScript+"\n"+psScript)
	return err
}
