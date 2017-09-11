package registry

import (
	"log"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
)

func PruneTimeTags(svcName, env string, n int) {
	if svcName == `` {
		for _, svcName := range conf.ServiceNames(env) {
			pruneTimeTags(svcName, env, n)
		}
	} else {
		pruneTimeTags(svcName, env, n)
	}
}

func pruneTimeTags(svcName, env string, n int) {
	tags := getTimeTagsSlice(svcName, env)
	if len(tags) <= n {
		return
	}
	RemoveTimeTags(svcName, env, tags[n:])
}

func RemoveTimeTags(imgName, env string, tags []string) {
	for _, tag := range tags {
		Remove(imgName, Digest(imgName, env+tag))
		log.Printf("removed %s:%s\n", imgName, env+tag)
	}
}
