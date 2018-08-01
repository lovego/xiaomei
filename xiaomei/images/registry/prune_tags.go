package registry

import (
	"fmt"
	"log"
	"strings"

	"github.com/lovego/cmd"
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
	reservedTags := reservedEnvTags(env, n, envTagsMap)
	imgName := conf.GetService(svcName, env).ImageName()
	toRemove := uniqDigestByTags(imgName, env, curEnvTags[n:])
	reserved := reservedDigest(svcName, env, reservedTags)
	for digest := range toRemove {
		if reserved[digest] {
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
		removeLocalImagesByTags(imgName, envTags)
	}
}

func removeLocalImagesByTags(imgName string, tags []string) {
	if len(tags) == 0 {
		return
	}
	images := []string{}
	for _, tag := range tags {
		images = append(images, fmt.Sprintf(" %s:%s", imgName, tag))
	}
	cmd.Run(cmd.O{}, `bash`, `-c`, `docker image rm -f`+strings.Join(images, ``))
	log.Printf("removed %s", strings.Join(images, "\n"))
}

// 当前环境前n个tag、其他环境前10个tag reserved
func reservedEnvTags(env string, n int, envTagsMap map[string][]string) map[string][]string {
	reservedTags := make(map[string][]string)
	for tagEnv, envTags := range envTagsMap {
		if tagEnv == env {
			reservedTags[tagEnv] = envTags[:n]
		} else if len(envTags) > 10 {
			reservedTags[tagEnv] = envTags[:10]
		}
	}
	return reservedTags
}

func reservedDigest(svcName, curEnv string, envTagsMap map[string][]string) map[string]bool {
	envDigest := make(map[string]bool)
	for tagEnv, envTags := range envTagsMap {
		if tagEnv != curEnv && len(envTags) > 10 {
			envTags = envTags[:10]
		}
		imgName := conf.GetService(svcName, tagEnv).ImageName()
		digestTagsMap := uniqDigestByTags(imgName, tagEnv, envTags)
		for digest := range digestTagsMap {
			envDigest[digest] = true
		}
	}
	return envDigest
}

func uniqDigestByTags(imgName, env string, tags []string) map[string][]string {
	digestMap := make(map[string][]string)
	for _, tag := range tags {
		envTag := env + `-` + tag
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
