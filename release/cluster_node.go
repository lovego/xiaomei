package release

import (
	"log"
	"os"
	"os/signal"
	"os/user"
	"strings"

	"github.com/lovego/addrs"
	"github.com/lovego/cmd"
)

type Node struct {
	user string
	// addr for ssh and service may be on different network segments.
	Addr        string            `yaml:"addr"`       // addr for ssh
	ServiceAddr string            `yaml:"seviceAddr"` // addr for service, use Addr if empty.
	Labels      map[string]string `yaml:"labels"`
}

func (n Node) Services(env, svcName string) []string {
	var svcNames []string
	if svcName == `` {
		svcNames = ServiceNames(env)
	} else {
		svcNames = []string{svcName}
	}
	var svcs []string
	for _, svcName := range svcNames {
		service := GetService(env, svcName)
		if n.Match(service.Nodes) {
			svcs = append(svcs, svcName)
		}
	}
	return svcs
}

func (n Node) Match(labels map[string]string) bool {
	for key, value := range labels {
		if nodeValue, ok := n.Labels[key]; !ok || nodeValue != value {
			return false
		}
	}
	return true
}

func (n Node) SshAddr() string {
	return strings.Split(n.user, `,`)[0] + `@` + n.Addr
}

func (n Node) SshCmd() string {
	return `ssh -t ` + n.SshAddr()
}

func (n Node) GetServiceAddr() string {
	if n.ServiceAddr != `` {
		return n.ServiceAddr
	}
	return n.Addr
}

func (n Node) Run(o cmd.O, script string) (string, error) {
	isLocal, err := n.IsLocalHost()
	if err != nil {
		return ``, err
	}
	if isLocal {
		c := make(chan os.Signal, 1)
		if o.Stdin == nil {
			signal.Notify(c, os.Interrupt)
			defer signal.Stop(c)
		}
		return cmd.Run(o, `bash`, `-c`, script)
	} else {
		return cmd.SshRun(o, n.SshAddr(), script)
	}
}

func (n Node) IsLocalHostP() bool {
	ok, err := n.IsLocalHost()
	if err != nil {
		log.Panic(err)
	}
	return ok
}

func (n Node) IsLocalHost() (bool, error) {
	return addrs.IsLocalhost(n.Addr)
}

func IsCurrentUser(users string) (bool, error) {
	u, err := user.Current()
	if u == nil || err != nil {
		return false, err
	}
	for _, user := range strings.Split(users, `,`) {
		if u.Username == user || u.Name == user {
			return true, nil
		}
	}
	return false, nil
}
