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
	removeTimeTags(svcName, env, tags[n:])
}

func RemoveTimeTags(svcName, env string, tags []string) {
	if svcName == `` {
		for _, svcName := range conf.ServiceNames(env) {
			removeTimeTags(svcName, env, tags)
		}
	} else {
		removeTimeTags(svcName, env, tags)
	}
}

func removeTimeTags(svcName, env string, tags []string) {
	imgName := conf.GetService(svcName, env).ImageName()
	for _, tag := range tags {
		if digest := Digest(imgName, env+tag); digest != `` {
			// TODO
			// it actually delete the image identified by the digest.
			// so the other tags on the same image are also deleted.
			// we want to just delete the tag. but we didn't see a remove tag registry api.
			Remove(imgName, digest)
			log.Printf("removed %s:%s\n", imgName, env+tag)
		}
	}
}
