package host

import (
	"fmt"
)

func (d driver) FlagsForRun(svcName string) ([]string, error) {
	return []string{
		`--network=host`,
		fmt.Sprintf(`-e=%s=%s`, portEnvName(svcName), portsOf(svcName)[0]),
	}, nil
}

func portEnvName(svcName string) string {
	switch svcName {
	case `app`:
		return `GOPORT`
	case `web`, `access`:
		return `NGPORT`
	default:
		panic(`unexpected svcName: ` + svcName)
	}
}
