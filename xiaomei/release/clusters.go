package release

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var theClusters map[string]Cluster

func GetClusters() map[string]Cluster {
	if theClusters == nil {
		content, err := ioutil.ReadFile(filepath.Join(Root(), `clusters.yml`))
		if err != nil {
			panic(err)
		}
		clusters := make(map[string]Cluster)
		if err := yaml.Unmarshal(content, clusters); err != nil {
			panic(err)
		}
		for _, cluster := range clusters {
			cluster.setNodesUser()
		}
		theClusters = clusters
	}
	return theClusters
}
