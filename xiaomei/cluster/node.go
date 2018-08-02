package cluster

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
	user       string
	jumpAddr   string
	Addr       string            `yaml:"addr"`
	Labels     map[string]string `yaml:"labels"`
	ListenAddr string            `yaml:"listenAddr"` // only for manager
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
	cmd := `ssh -t ` + n.SshAddr()
	if n.jumpAddr != `` {
		cmd = `ssh -t ` + n.jumpAddr + ` ` + cmd
	}
	return cmd
}

func (n Node) GetListenAddr() string {
	if n.ListenAddr != `` {
		return n.ListenAddr
	}
	return n.Addr
}

func (n Node) Run(o cmd.O, script string) (string, error) {
	isLocal, err := IsLocalHost(n.Addr)
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
		if n.jumpAddr == `` {
			return cmd.SshRun(o, n.SshAddr(), script)
		} else {
			return cmd.SshJumpRun(o, n.jumpAddr, n.SshAddr(), script)
		}
	}
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

func IsLocalHost(addr string) (bool, error) {
	ips, err := net.LookupIP(addr)
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
