package config

import (
	"net"
	"strings"

	"github.com/bughou-go/xiaomei/utils/slice"
)

type Server struct {
	Addr       string   `yaml:"addr"`
	ListenAddr string   `yaml:"listenAddr"`
	Tasks      []string `yaml:"tasks"`
	AppStartOn string   `yaml:"appStartOn"`
}

func (s *Server) IsLocal() bool {
	return slice.ContainsString(MachineAddrs(), s.Addr)
}

func (s *Server) HasTask(name string) bool {
	if name == `` {
		return true
	}
	return slice.ContainsString(s.Tasks, name)
}

func (s *Server) Match(feature string) bool {
	if feature == `` {
		return true
	}
	return slice.ContainsString(s.Tasks, feature) ||
		strings.Contains(s.Addr, feature) ||
		strings.Contains(s.ListenAddr, feature)
}

func (s *Server) AppAddr() string {
	if s.ListenAddr != `` {
		return s.ListenAddr + `:` + App.Port()
	}
	return s.Addr + `:` + App.Port()
}

func (s *Server) GodocAddr() string {
	if s.ListenAddr != `` {
		return s.ListenAddr + `:` + Godoc.Port()
	}
	return s.Addr + `:` + Godoc.Port()
}

func (s *Server) SshAddr() string {
	return Deploy.User() + `@` + s.Addr
}

var machineAddrs []string

func MachineAddrs() []string {
	if machineAddrs != nil {
		return machineAddrs
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
	machineAddrs = addrs
	return machineAddrs
}
