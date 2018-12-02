package cluster

import (
	"io/ioutil"
	"path/filepath"

	"github.com/lovego/xiaomei/release"
	"gopkg.in/yaml.v2"
)

func init() {
	release.RegisterEnvsGetter(Envs)
}

var theClusters map[string]*Cluster

func GetClusters() map[string]*Cluster {
	if theClusters == nil {
		content, err := ioutil.ReadFile(filepath.Join(release.Root(), `clusters.yml`))
		if err != nil {
			panic(err)
		}
		clusters := make(map[string]*Cluster)
		if err := yaml.Unmarshal(content, clusters); err != nil {
			panic(err)
		}
		for env, cluster := range clusters {
			cluster.init(env)
		}
		theClusters = clusters
	}
	return theClusters
}

var theEnvs []string

func Envs() []string {
	if theEnvs == nil {
		envs := []string{}
		// if release.InProject() {
		for env := range GetClusters() {
			envs = append(envs, env)
		}
		// }
		theEnvs = envs
	}
	return theEnvs
}
