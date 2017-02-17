package config

import (
	"path/filepath"
)

var Cluster ClusterConf

type ClusterConf struct {
	conf clusterConf
}

type clusterConf struct {
	User       string `yaml:"user"`
	DeployRoot string `yaml:"deployRoot"`
	Managers   []Node `yaml:"managers"`
	Workers    []Node `yaml:"workers"`
}

type Node struct {
	Addr       string   `yaml:"addr"`
	Labels     []string `yaml:"labels"`
	ListenAddr string   `yaml:"listenAddr"` // only for manager
}

func (c *ClusterConf) User() string {
	Load()
	return c.conf.User
}

func (c *ClusterConf) DeployName() string {
	Load()
	return App.Name() + `_` + App.Env()
}

func (c *ClusterConf) DeployRoot() string {
	Load()
	return c.conf.DeployRoot
}

func (c *ClusterConf) DeployPath() string {
	Load()
	return filepath.Join(c.DeployRoot(), c.DeployName())
}

func (c *ClusterConf) Managers() []Node {
	Load()
	return c.conf.Managers
}

func (c *ClusterConf) Workers() []Node {
	Load()
	return c.conf.Workers
}

func (n Node) SshAddr() string {
	return Cluster.User() + `@` + n.Addr
}
