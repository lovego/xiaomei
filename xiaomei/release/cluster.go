package release

func GetCluster() Cluster {
	cluster, ok := GetClusters()[Env()]
	if !ok {
		panic(`empty cluster config for env: ` + Env())
	}
	return cluster
}

type Cluster struct {
	User     string `yaml:"user"`
	Managers []Node `yaml:"managers"`
	Workers  []Node `yaml:"workers"`
}

func (c *Cluster) setNodesUser() {
	for i := range c.Managers {
		c.Managers[i].user = c.User
	}
	for i := range c.Workers {
		c.Workers[i].user = c.User
	}
}

func (c Cluster) SshAddr() string {
	if len(c.Managers) == 0 {
		panic(`the cluster have no managers.`)
	}
	m := c.Managers[0]
	return m.SshAddr()
}

func (c Cluster) List() ([]string, []string) {
	ms := []string{}
	for _, m := range c.Managers {
		ms = append(ms, m.SshAddr())
	}
	ws := []string{}
	for _, w := range c.Workers {
		ws = append(ws, w.SshAddr())
	}
	return ms, ws
}
