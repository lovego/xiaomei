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
	envTagsMap := getEnvTimeTagsMap(svcName, env)
	curEnvTags := envTagsMap[env]
	if len(curEnvTags) <= n {
		return
	}
	delete(envTagsMap, env)
	imgName := conf.GetService(svcName, env).ImageName()
	reserved := uniqDigestByTags(imgName, env, curEnvTags[:n])
	toRemove := uniqDigestByTags(imgName, env, curEnvTags[n:])
	filterNotRemove(svcName, toRemove, reserved, envTagsMap)
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

func filterNotRemove(svcName string, toRemove, reserved, otherEnvTags map[string][]string) {
	for digest, _ := range toRemove {
		if reserved[digest] != nil {
			delete(toRemove, digest)
			continue
		}
		for tagEnv, envTags := range otherEnvTags {
			if len(envTags) > 10 {
				envTags = envTags[:10]
			}
			// 默认判断其他环境前10个tag
			imgName := conf.GetService(svcName, tagEnv).ImageName()
			top10Digests := uniqDigestByTags(imgName, tagEnv, envTags)
			if top10Digests[digest] != nil {
				delete(toRemove, digest)
			}
		}
	}
}
