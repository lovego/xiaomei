package cluster

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type ClusterConf struct {
	User     string     `yaml:"user"`
	Managers []NodeConf `yaml:"managers"`
	Workers  []NodeConf `yaml:"workers"`
}

type NodeConf struct {
	user       string
	Addr       string   `yaml:"addr"`
	Labels     []string `yaml:"labels"`
	ListenAddr string   `yaml:"listenAddr"` // only for manager
}

var theClusters map[string]ClusterConf

func GetConfig(env string) (ClusterConf, error) {
	if theClusters == nil {
		if clusters, err := loadClustersConfig(); err != nil {
			return ClusterConf{}, err
		} else {
			theClusters = clusters
		}
	}
	if env == `` {
		env = `dev`
	}
	if clusterConf, ok := theClusters[env]; ok {
		clusterConf.init()
		return clusterConf, nil
	}
	return ClusterConf{}, errors.New(`empty cluster config for env: ` + env)
}

func loadClustersConfig() (map[string]ClusterConf, error) {
	content, err := ioutil.ReadFile(filepath.Join(config.Root(), `../clusters.yml`))
	if err != nil {
		return nil, err
	}
	clusters := make(map[string]ClusterConf)
	if err := yaml.Unmarshal(content, clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func (c *ClusterConf) init() {
	for i := range c.Managers {
		c.Managers[i].user = c.User
	}
	for i := range c.Workers {
		c.Workers[i].user = c.User
	}
}

func (c ClusterConf) SshAddr() (string, error) {
	if len(c.Managers) == 0 {
		return ``, errors.New(`the cluster have no managers.`)
	}
	m := c.Managers[0]
	return c.User + `@` + m.Addr, nil
}

func (n NodeConf) SshAddr() string {
	return n.user + `@` + n.Addr
}
