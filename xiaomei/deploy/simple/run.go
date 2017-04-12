package simple

import (
	"fmt"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
)

func (d driver) FlagsForRun(svcName string) ([]string, error) {
	flags := []string{`--network=host`}
	portEnv := portEnvName(svcName)
	ports := simpleconf.PortsOf(svcName)
	if portEnv != `` && len(ports) > 0 {
		flags = append(flags, fmt.Sprintf(`-e=%s=%s`, portEnv, ports[0]))
	}
	return flags, nil
}

func portEnvName(svcName string) string {
	switch svcName {
	case `app`:
		return `GOPORT`
	case `web`, `access`:
		return `NGPORT`
	default:
		return ``
	}
}
