package cluster

import (
	"log"
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
)

type Cluster struct {
	env      string
	User     string `yaml:"user"`
	JumpAddr string `yaml:"jumpAddr"`
	Nodes    []Node `yaml:"nodes"`
}

func Get(env string) Cluster {
	cluster, ok := GetClusters()[env]
	if !ok {
		log.Fatalf("empty cluster config for env: %v", env)
	}
	return cluster
}

func (c *Cluster) init(env string) {
	c.env = env
	for i := range c.Nodes {
		c.Nodes[i].user = c.User
		c.Nodes[i].jumpAddr = c.JumpAddr
	}
}

func (c Cluster) GetNodes(feature string) (nodes []Node) {
	for _, node := range c.Nodes {
		if feature == `` || strings.Contains(node.Addr, feature) {
			nodes = append(nodes, node)
		}
	}
	return
}

func (c Cluster) Run(feature string, o cmd.O, script string) (string, error) {
	for _, node := range c.GetNodes(feature) {
		return node.Run(o, script)
	}
	return ``, nil
}

func (c Cluster) ServiceRun(svcName, feature string, o cmd.O, script string) (string, error) {
	labels := conf.GetService(svcName, c.env).Nodes
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
