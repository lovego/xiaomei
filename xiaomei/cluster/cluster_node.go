package cluster

type Node struct {
	user       string
	Addr       string   `yaml:"addr"`
	Labels     []string `yaml:"labels"`
	ListenAddr string   `yaml:"listenAddr"` // only for manager
}

func (n Node) SshAddr() string {
	return n.user + `@` + n.Addr
}
