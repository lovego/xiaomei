package host

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
)

func (d driver) Ps(svcName string, watch bool, options []string) error {
	var script string

	_, err := cluster.Run(cmd.O{}, script)
	return err
}
