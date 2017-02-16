package config

import (
	"path/filepath"
)

var Cluster ClusterConf

type ClusterConf struct {
	conf clusterConf
}

type clusterConf struct {
	User       string        `yaml:"user"`
	DeployRoot string        `yaml:"deployRoot"`
	Managers   []ManagerNode `yaml:"managers"`
	Workers    []WorkerNode  `yaml:"workers"`
}

type ManagerNode struct {
	Addr       string   `yaml:"addr"`
	ListenAddr string   `yaml:"listenAddr"`
	Labels     []string `yaml:"labels"`
}

type WorkerNode struct {
	Addr   string   `yaml:"addr"`
	Labels []string `yaml:"labels"`
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

func (c *ClusterConf) Managers() []ManagerNode {
	Load()
	return c.conf.Managers
}

func (c *ClusterConf) Workers() []WorkerNode {
	Load()
	return c.conf.Workers
}

type Node interface {
	SshAddr() string
}

func (n ManagerNode) SshAddr() string {
	return Cluster.User() + `@` + n.Addr
}

func (n WorkerNode) SshAddr() string {
	return Cluster.User() + `@` + n.Addr
}
