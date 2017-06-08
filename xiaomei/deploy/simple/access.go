package simple

import (
	"errors"
)

import (
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
)

func (d driver) Addrs(svcName string) (addrs []string, err error) {
	ports := simpleconf.PortsOf(svcName)
	for _, node := range nodesFor(svcName) {
		for _, port := range ports {
			addrs = append(addrs, node.GetListenAddr()+`:`+port)
		}
	}
	if len(addrs) == 0 {
		err = errors.New(`no instance defined for: ` + svcName)
	}
	return
}

func nodesFor(svcName string) (nodes []cluster.Node) {
	service := simpleconf.GetService(svcName)
	for _, node := range cluster.Nodes() {
		if node.Match(service.Nodes) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
