package config

import (
	"path/filepath"
)

var Cluster ClusterConf

type ClusterConf struct {
	conf clusterConf
}

type clusterConf struct {
	User       string   `yaml:"user"`
	DeployRoot string   `yaml:"deployRoot"`
	Masters    []string `yaml:"masters"`
	Workers    []string `yaml:"workers"`
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

func (c *ClusterConf) Masters() []string {
	Load()
	return c.conf.Masters
}

func (c *ClusterConf) Workers() []string {
	Load()
	return c.conf.Workers
}
