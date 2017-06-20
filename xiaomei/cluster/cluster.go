package cluster

import (
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func Run(o cmd.O, script string) (string, error) {
	return GetCluster().Manager().Run(o, script)
}

func Nodes() []Node {
	return GetCluster().Nodes()
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
	JumpAddr string `yaml:"jumpAddr"`
	Managers []Node `yaml:"managers"`
	Workers  []Node `yaml:"workers"`
	nodes    []Node
}

func (c *Cluster) init() {
	for i := range c.Managers {
		c.Managers[i].user = c.User
		c.Managers[i].jumpAddr = c.JumpAddr
	}
	for i := range c.Workers {
		c.Workers[i].user = c.User
		c.Workers[i].jumpAddr = c.JumpAddr
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
		ms = append(ms, m.SshCmd())
	}
	println("managers:\n", strings.Join(ms, "\n"))
	ws := []string{}
	for _, w := range c.Workers {
		ws = append(ws, w.SshCmd())
	}
	println("workers:\n", strings.Join(ws, "\n"))
}
