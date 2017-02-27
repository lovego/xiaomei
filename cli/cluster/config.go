package config

var Cluster ClusterConf

type ClusterConf struct {
	conf clusterConf
}

type clusterConf struct {
	User     string `yaml:"user"`
	Managers []Node `yaml:"managers"`
	Workers  []Node `yaml:"workers"`
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
