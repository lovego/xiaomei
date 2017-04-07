package host

import (
	"fmt"
	"regexp"
	"strconv"
)

var Driver driver

type driver struct{}

func (d driver) ImageNameOf(svcName string) string {
	svc := getService(svcName)
	if svc.Image == `` {
		panic(fmt.Sprintf(`release.yml: %s.image: empty.`, svcName))
	}
	return svc.Image
}

var rePort = regexp.MustCompile(`^\d+$`)
var rePorts = regexp.MustCompile(`^(\d+)-(\d+)$`)

func (d driver) PortsOf(svcName string) (ports []string) {
	svc := getService(svcName)
	if svc.Ports == `` {
		return
	}
	if rePort.MatchString(svc.Ports) {
		ports = append(ports, svc.Ports)
	} else if m := rePorts.FindStringSubmatch(svc.Ports); len(m) == 3 {
		start, err := strconv.Atoi(m[1])
		if err != nil {
			panic(err)
		}
		end, err := strconv.Atoi(m[2])
		if err != nil {
			panic(err)
		}
		for ; start <= end; start++ {
			ports = append(ports, strconv.Itoa(start))
		}
	} else {
		panic(fmt.Sprintf(`release.yml: %s.ports: illegal format.`, svcName))
	}

	return
}

func (d driver) ServiceNames() map[string]bool {
	m := make(map[string]bool)
	for svcName := range getRelease() {
		m[svcName] = true
	}
	return m
}
