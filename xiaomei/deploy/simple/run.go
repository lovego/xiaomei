package simple

import (
	"fmt"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func (d driver) FlagsForRun(svcName string) ([]string, error) {
	flags := []string{`--network=host`, `--name=` + release.DeployName() + `_` + svcName + `.run`}
	if portEnvName := images.Get(svcName).PortEnvName(); portEnvName != `` {
		if ports := simpleconf.PortsOf(svcName); len(ports) > 0 {
			flags = append(flags, fmt.Sprintf(`-e=%s=%s`, portEnvName, ports[0]))
		}
	}
	return flags, nil
}
