package config

import "strings"

type Server struct {
	Addr       string   `yaml:"addr"`
	ListenAddr string   `yaml:"listenAddr"`
	Tasks      []string `yaml:"tasks"`
	AppStartOn string   `yaml:"appStartOn"`
}

func (s *Server) HasTask(name string) bool {
	return contains(s.Tasks, name)
}

func (s *Server) HasAddr(addrs []string) bool {
	return contains(addrs, s.Addr) || contains(addrs, s.ListenAddr)
}

func (s *Server) Match(feature string) bool {
	if feature == `` {
		return true
	}
	return contains(s.Tasks, feature) ||
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
