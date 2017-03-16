package cluster

import (
	"github.com/bughou-go/xiaomei/utils/cmd"
)

type Node struct {
	user       string
	Addr       string   `yaml:"addr"`
	Labels     []string `yaml:"labels"`
	ListenAddr string   `yaml:"listenAddr"` // only for manager
}

func (n Node) SshAddr() string {
	return n.user + `@` + n.Addr
}

func (n Node) Run(o cmd.O, script string) (string, error) {
	return cmd.SshRun(o, n.SshAddr(), script)
}
