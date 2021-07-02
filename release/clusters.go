package release

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/lovego/cmd"
	yaml "gopkg.in/yaml.v2"
)

var theClusters map[string]*Cluster

type Cluster struct {
	env   string
	User  string `yaml:"user"`
	Nodes []Node `yaml:"nodes"`
}

func GetClusters() map[string]*Cluster {
	if theClusters == nil {
		content, err := ioutil.ReadFile(filepath.Join(Root(), `clusters.yml`))
		if err != nil {
			panic(err)
		}
		theClusters = make(map[string]*Cluster)
		if err := yaml.Unmarshal(content, theClusters); err != nil {
			panic(err)
		}
		for env, cluster := range theClusters {
			cluster.init(env)
		}
	}
	return theClusters
}

func GetCluster(env string) *Cluster {
	cluster := GetClusters()[env]
	if cluster == nil {
		log.Fatalf("empty cluster config for env: %v", env)
	}
	return cluster
}

func (c *Cluster) init(env string) {
	c.env = env
	for i := range c.Nodes {
		c.Nodes[i].user = c.User
	}
}

func (c Cluster) IsLocalHost() (bool, error) {
	for _, node := range c.Nodes {
		if ok, err := node.IsLocalHost(); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (c Cluster) GetNodes(feature string) (nodes []Node) {
	for _, node := range c.Nodes {
		if feature == `` || strings.Contains(node.Addr, feature) {
			nodes = append(nodes, node)
		}
	}
	return
}

func (c Cluster) NodesCount() int {
	return len(c.Nodes)
}

func (c Cluster) Run(feature string, o cmd.O, script string) (string, error) {
	for _, node := range c.GetNodes(feature) {
		if node.IsLocalHostP() {
			return node.Run(o, script)
		}
	}
	for _, node := range c.GetNodes(feature) {
		return node.Run(o, script)
	}
	return ``, nil
}

func (c Cluster) ServiceRun(svcName, feature string, o cmd.O, script string) (string, error) {
	labels := GetService(c.env, svcName).Nodes
	for _, node := range c.GetNodes(feature) {
		if node.IsLocalHostP() && node.Match(labels) {
			return node.Run(o, script)
		}
	}
	for _, node := range c.GetNodes(feature) {
		if node.Match(labels) {
			return node.Run(o, script)
		}
	}
	return ``, nil
}

func (c Cluster) List() {
	nodes := []string{}
	for _, m := range c.Nodes {
		nodes = append(nodes, m.SshCmd())
	}
	println(strings.Join(nodes, "\n"))
}
