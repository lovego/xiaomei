package simple

import (
	"fmt"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func (d driver) FlagsForRun(svcName string) ([]string, error) {
	name := release.Name() + `_` + svcName
	portEnv := ``
	if portEnvName := images.Get(svcName).PortEnvName(); portEnvName != `` {
		if ports := simpleconf.PortsOf(svcName); len(ports) > 0 {
			name += `.` + ports[0]
			portEnv = fmt.Sprintf(`-e=%s=%s`, portEnvName, ports[0])
		}
	}
	flags := []string{`--network=host`, `--name=` + name}
	if portEnv != `` {
		flags = append(flags, portEnv)
	}
	return flags, nil
}
