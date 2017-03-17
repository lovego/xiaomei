package cluster

import (
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/slice"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Run(o cmd.O, script string) (string, error) {
	return GetCluster().Manager().Run(o, script)
}

func Nodes() []Node {
	return GetCluster().Nodes()
}

func AccessNodes() (result []Node) {
	for _, node := range GetCluster().Nodes() {
		if slice.ContainsString(node.Labels, `hasAccess=true`) {
			result = append(result, node)
		}
	}
	return
}

func GetCluster() Cluster {
	cluster, ok := GetClusters()[release.Env()]
	if !ok {
		panic(`empty cluster config for env: ` + release.Env())
	}
	return cluster
}

type Cluster struct {
	User     string `yaml:"user"`
	Managers []Node `yaml:"managers"`
	Workers  []Node `yaml:"workers"`
	nodes    []Node
}

func (c *Cluster) setNodesUser() {
	for i := range c.Managers {
		c.Managers[i].user = c.User
	}
	for i := range c.Workers {
		c.Workers[i].user = c.User
	}
}

func (c Cluster) Manager() Node {
	if len(c.Managers) == 0 {
		panic(`the cluster have no managers.`)
	}
	return c.Managers[0]
}

func (c Cluster) Nodes() []Node {
	if c.nodes == nil {
		c.nodes = append(c.nodes, c.Managers...)
		c.nodes = append(c.nodes, c.Workers...)
	}
	return c.nodes
}

func (c Cluster) List() {
	ms := []string{}
	for _, m := range c.Managers {
		ms = append(ms, m.SshAddr())
	}
	println(`managers: `, strings.Join(ms, " \t"))
	ws := []string{}
	for _, w := range c.Workers {
		ws = append(ws, w.SshAddr())
	}
	println(`workers: `, strings.Join(ws, " \t"))
}
