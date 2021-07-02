package release

import (
	"net"
	"os"
	"os/signal"
	"os/user"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/slice"
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
		panic(err)
	}
	return ok
}

func (n Node) IsLocalHost() (bool, error) {
	ips, err := net.LookupIP(n.Addr)
	if err != nil {
		return false, err
	}
	for _, ip := range ips {
		if ip.IsLoopback() {
			return true, nil
		}
	}
	machineAddrs := MachineAddrs()
	for _, ip := range ips {
		if slice.ContainsString(machineAddrs, ip.String()) {
			return true, nil
		}
	}
	return false, nil
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

var theMachineAddrs []string

func MachineAddrs() []string {
	if theMachineAddrs != nil {
		return theMachineAddrs
	}

	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	addrs := make([]string, len(ifcAddrs))
	for i, ifcAddr := range ifcAddrs {
		addr := ifcAddr.String()
		if i := strings.IndexByte(addr, '/'); i >= 0 {
			addr = addr[:i]
		}
		addrs[i] = addr
	}
	theMachineAddrs = addrs
	return theMachineAddrs
}
