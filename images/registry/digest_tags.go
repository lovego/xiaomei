package registry

import (
	"fmt"

	"github.com/lovego/xiaomei/deploy/conf"
)

func DigestTimeTags(svcName, env string, tags []string) {
	if svcName == `` {
		for _, svcName := range conf.ServiceNames(env) {
			digestTimeTags(svcName, env, tags)
		}
	} else {
		digestTimeTags(svcName, env, tags)
	}
}

func digestTimeTags(svcName, env string, tags []string) {
	imgName := conf.GetService(svcName, env).ImageName()
	for _, tag := range tags {
		if digest := Digest(imgName, env+tag); digest != `` {
			fmt.Printf("%s:%s %s\n", imgName, env+tag, Digest(imgName, digest))
		}
	}
}
