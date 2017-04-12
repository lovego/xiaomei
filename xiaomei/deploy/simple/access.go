package simple

import (
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
)

func (d driver) AccessAddrs(svcName string) (addrs []string) {
	ports := simpleconf.portsOf(svcName)
	for _, node := range cluster.Nodes() {
		for _, port := range ports {
			addrs = append(addrs, node.GetListenAddr()+`:`+port)
		}
	}
	return
}
