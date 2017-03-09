package cluster

import (
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Run(o cmd.O, script string) error {
	_, err := cmd.SshRun(o, release.GetCluster().SshAddr(), script)
	return err
}
