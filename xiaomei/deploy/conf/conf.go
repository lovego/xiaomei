package conf

import (
	"regexp"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	"github.com/lovego/xiaomei/xiaomei/release"
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

var reImageName = regexp.MustCompile(`^(.*):([\w.-]+)$`)

func ImageNameOf(svcName string) string {
	name := simpleconf.ImageNameOf(svcName)
	if reImageName.MatchString(name) {
		name += `:` + release.Env()
	}
	return name
}

func ImageNameAndTagOf(svcName string) (name, tag string) {
	name = simpleconf.ImageNameOf(svcName)
	if m := reImageName.FindStringSubmatch(name); len(m) == 3 {
		return m[1], m[2]
	} else {
		return name, release.Env()
	}
}

func VolumesFor(svcName string) []string {
	return simpleconf.VolumesFor(svcName)
}

func CommandFor(svcName string) []string {
	return simpleconf.CommandFor(svcName)
}
