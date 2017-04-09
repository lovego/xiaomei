package host

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/lovego/xiaomei/xiaomei/cluster"
)

func (d driver) AccessAddrs(svcName string) (addrs []string) {
	ports := portsOf(svcName)
	for _, node := range cluster.Nodes() {
		for _, port := range ports {
			addrs = append(addrs, node.GetListenAddr()+`:`+port)
		}
	}
	return
}

var rePort = regexp.MustCompile(`^\d+$`)
var rePorts = regexp.MustCompile(`^(\d+)-(\d+)$`)

func portsOf(svcName string) (ports []string) {
	svc := getService(svcName)
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
	if len(ports) == 0 {
		panic(fmt.Sprintf(`release.yml: %s.ports: can't be empty.`, svcName))
	}
	return
}
