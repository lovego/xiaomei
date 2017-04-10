package host

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func (d driver) Ps(svcName string, watch bool, options []string) error {
	script := fmt.Sprintf(`docker ps -f name=%s_%s`, release.Name(), svcName)
	if watch {
		script = `watch ` + script
	}
	_, err := cluster.Run(cmd.O{}, script)
	return err
}
