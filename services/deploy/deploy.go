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
	svcName, env, timeTag, filter       string
	noPullBaseImage, noPushImageIfLocal bool
	beforeScript                        string
	noBeforeScriptOnLocal, noWatch      bool
}

func (d Deploy) start() error {
	if d.beforeScript != `` && !d.noBeforeScriptOnLocal {
		if _, err := cmd.Run(cmd.O{}, `bash`, `-c`, d.beforeScript); err != nil {
			return err
		}
	}
	if d.timeTag == `` {
		d.timeTag = release.TimeTag(d.env)
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
	if err := images.Build(d.svcName, d.env, d.timeTag, !d.noPullBaseImage); err != nil {
		return err
	}
	if d.noPushImageIfLocal {
		if ok, err := release.GetCluster(d.env).IsLocalHost(); err != nil {
			return err
		} else if ok {
			return nil
		}
	}
	return images.Push(d.svcName, d.env, d.timeTag)
}

func (d Deploy) run() error {
	psScript := fmt.Sprintf(` docker ps -f name=^/%s`, release.ServiceName(d.svcName, d.env))
	if !d.noWatch {
		psScript = oam.WatchCmd() + psScript
	}
	expectHighAvailable := len(release.GetCluster(d.env).GetNodes("")) >= 2
	var recoverAccess bool
	for _, node := range release.GetCluster(d.env).GetNodes(d.filter) {
		if svcs := node.Services(d.env, d.svcName); len(svcs) > 0 {
			if access.HasAccess(svcs) {
				if expectHighAvailable {
					if err := access.SetupNginx(d.env, "", node.Addr); err != nil {
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
		if expectHighAvailable {
			return access.SetupNginx(d.env, "", "")
		} else {
			return access.ReloadNginx(d.env, "")
		}
	}
	return nil
}

func (d Deploy) runNode(svcs []string, node release.Node, psScript string) error {
	log.Println(color.GreenString(`deploying ` + node.SshAddr()))
	deployScript, err := getDeployScript(svcs, d.env, d.timeTag)
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
