package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Deploy(env, svc string) error {
	stackContent, err := GetDeployStackContent(svc)
	if err != nil {
		return err
	}
	clusterConf, err := cluster.GetConfig(env)
	if err != nil {
		return err
	}
	addr, err := clusterConf.SshAddr()
	if err != nil {
		return err
	}
	stackFile = config.DeployName + `_`
	cmd.SshRun(cmd.O{}, addr, `cat - > %s; docker stack deploy --compose-file=%s`)
}

func GetDeployStackContent(svc string) (string, error) {
	var stackContent string
	if svc == `` {
		if _stack, err := GetStackFile(); err == nil {
			stack = _stack
		} else {
			return ``, err
		}
	} else {
	}
}
