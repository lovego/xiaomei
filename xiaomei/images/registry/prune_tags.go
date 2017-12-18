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
	imgName := conf.GetService(svcName, env).ImageName()
	reserved := uniqDigestByTags(imgName, env, tags[:n])
	toRemove := uniqDigestByTags(imgName, env, tags[n:])
	for digest, _ := range toRemove {
		if reserved[digest] == nil {
			delete(toRemove, digest)
		}
	}
	removeTimeTags(imgName, toRemove)
}

func RemoveTimeTags(svcName, env string, tags []string) {
	if svcName == `` {
		for _, svcName := range conf.ServiceNames(env) {
			imgName := conf.GetService(svcName, env).ImageName()
			removeTimeTags(imgName, uniqDigestByTags(imgName, env, tags))
		}
	} else {
		imgName := conf.GetService(svcName, env).ImageName()
		removeTimeTags(imgName, uniqDigestByTags(imgName, env, tags))
	}
}

func removeTimeTags(imgName string, toRemove map[string][]string) {
	for digest, envTags := range toRemove {
		Remove(imgName, digest)
		for _, envTag := range envTags {
			log.Printf("removed %s:%s\n", imgName, envTag)
		}
	}
}

func uniqDigestByTags(imgName, env string, tags []string) map[string][]string {
	digestMap := make(map[string][]string)
	for _, tag := range tags {
		envTag := env + tag
		digest := Digest(imgName, envTag)
		envTags := digestMap[digest]
		if envTags == nil {
			envTags = []string{envTag}
		} else {
			envTags = append(envTags, envTag)
		}
		digestMap[digest] = envTags
	}
	return digestMap
}
