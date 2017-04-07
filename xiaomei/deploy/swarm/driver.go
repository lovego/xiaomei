package swarm

import (
	"fmt"
	"strings"
)

var Driver driver

type driver struct{}

func (d driver) ImageNameOf(svcName string) string {
	svc := getService(svcName)
	image := svc[`image`]
	if image == nil {
		panic(fmt.Sprintf(`stack.yml: services.%s.image: undefined.`, svcName))
	}
	if str, ok := image.(string); ok && str != `` {
		return str
	} else {
		panic(fmt.Sprintf(`stack.yml: services.%s.image: should be a string.`, svcName))
	}
}

func (d driver) PortsOf(svcName string) (ports []string) {
	service := getService(svcName)
	iports, _ := service[`ports`].([]interface{})
	for i, iport := range iports {
		switch port := iport.(type) {
		case string:
			if index := strings.IndexByte(port, ':'); index < 0 {
				panic(fmt.Sprintf(
					`stack.yml: services.%s.ports: random host port is not allowed: %s.`, svcName, port,
				))
			}
			ports = append(ports, port)
		case map[interface{}]interface{}:
			ports = append(ports, fmt.Sprint(port[`published`])+`:`+fmt.Sprint(port[`target`]))
		default:
			panic(fmt.Sprintf(
				`stack.yml: services.%s.ports[%d]: should be a string, got: %s.`, svcName, i, iport,
			))
		}
	}
	return ports
}

func (d driver) ServiceNames() map[string]bool {
	m := make(map[string]bool)
	for svcName := range getStack().Services {
		m[svcName] = true
	}
	return m
}
