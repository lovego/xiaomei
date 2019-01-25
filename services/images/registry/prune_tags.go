package registry

import (
	"fmt"
	"log"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release/cluster"
	"github.com/lovego/xiaomei/services/deploy/conf"
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
	imgName := conf.GetService(svcName, env).ImageName()
	toRemove := getToRemoveTimeTags(imgName, svcName, env, n)
	removeTimeTags(imgName, env, toRemove)
}

func getToRemoveTimeTags(imgName, svcName, env string, n int) map[string][]string {
	envTagsMap := getEnvTimeTagsMap(svcName, env)
	curEnvTags := envTagsMap[env]
	if len(curEnvTags) <= n {
		return nil
	}
	reservedTags := reservedEnvTags(env, n, envTagsMap)
	toRemove := uniqDigestByTags(imgName, env, curEnvTags[n:])
	reserved := reservedDigest(svcName, env, reservedTags)
	for digest := range toRemove {
		if digest == `` {
			continue
		}
		if reserved[digest] {
			delete(toRemove, digest)
		}
	}
	return toRemove
}

func RemoveTimeTags(svcName, env string, tags []string) {
	if svcName == `` {
		for _, svcName := range conf.ServiceNames(env) {
			imgName := conf.GetService(svcName, env).ImageName()
			removeTimeTags(imgName, env, uniqDigestByTags(imgName, env, tags))
		}
	} else {
		imgName := conf.GetService(svcName, env).ImageName()
		removeTimeTags(imgName, env, uniqDigestByTags(imgName, env, tags))
	}
}

func removeTimeTags(imgName, env string, toRemove map[string][]string) {
	if len(toRemove) == 0 {
		return
	}
	toRemoveTags := []string{}
	for digest, envTags := range toRemove {
		if digest != `` {
			Remove(imgName, digest)
		}
		toRemoveTags = append(toRemoveTags, envTags...)
	}
	rmiScript := getRmiLocalScript(imgName, toRemoveTags)
	removeNodesTimeTags(env, rmiScript)
	cmd.Run(cmd.O{}, `bash`, `-c`, rmiScript)
	for _, tag := range toRemoveTags {
		log.Printf("removed %s:%s", imgName, tag)
	}
}

// 删除环境运行机器上本地的镜像
func removeNodesTimeTags(env, rmiScript string) {
	if rmiScript == `` {
		return
	}
	for _, node := range cluster.Get(env).GetNodes(``) {
		if isLocalHost, err := node.IsLocalHost(); err != nil {
			log.Println(err)
			continue
		} else if isLocalHost {
			continue
		} else {
			if _, err := node.Run(cmd.O{}, rmiScript); err != nil {
				log.Println(err)
			}
		}
	}
}

func getRmiLocalScript(imgName string, tags []string) string {
	if len(tags) == 0 {
		return ``
	}
	images := []string{}
	for _, tag := range tags {
		images = append(images, fmt.Sprintf(" %s:%s", imgName, tag))
	}
	return `set -e; docker image rm -f` + strings.Join(images, ``)
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
			if digest != `` {
				envDigest[digest] = true
			}
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
