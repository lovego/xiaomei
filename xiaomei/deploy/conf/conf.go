package conf

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	// "github.com/lovego/xiaomei/xiaomei/deploy/stackconf"
)

func Type() string {
	return `simple`
}

func File() string {
	return `simple.yml`
}

func ServiceNames() map[string]bool {
	return simpleconf.ServiceNames()
}

func ImageNameOf(svcName string) string {
	return simpleconf.ImageNameOf(svcName)
}

func VolumesFor(svcName string) []string {
	return simpleconf.VolumesFor(svcName)
}

func CommandFor(svcName string) []string {
	return simpleconf.CommandFor(svcName)
}
