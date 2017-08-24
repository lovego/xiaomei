package cluster

import (
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func Run(o cmd.O, feature, script string) (string, error) {
	for _, node := range Nodes(feature) {
		return node.Run(o, script)
	}
	return ``, nil
}

func ServiceRun(o cmd.O, svcName, feature, script string) (string, error) {
	labels := conf.GetService(svcName).Nodes
	for _, node := range Nodes(feature) {
		if node.Match(labels) {
			return node.Run(o, script)
		}
	}
	return ``, nil
}

func Nodes(feature string) []Node {
	return GetCluster().Nodes(feature)
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

func (c Cluster) Nodes(feature string) (nodes []Node) {
	for _, node := range c.Managers {
		if feature == `` || strings.Contains(node.Addr, feature) {
			nodes = append(nodes, node)
		}
	}
	for _, node := range c.Workers {
		if feature == `` || strings.Contains(node.Addr, feature) {
			nodes = append(nodes, node)
		}
	}
	return
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
