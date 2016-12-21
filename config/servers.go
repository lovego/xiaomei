package config

import (
	"net"
	"strings"
)

var Servers serverVar

type serverVar struct {
	conf ServersConf
}

type ServersConf []Server

type Server struct {
	Addr       string `yaml:"addr"`
	Tasks      string `yaml:"tasks"`
	AppAddr    string `yaml:"appAddr"`
	AppStartOn string `yaml:"appStartOn"`
}

func (s *serverVar) All() []Server {
	Load()
	return s.conf
}

func (s *serverVar) Current() Server {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range s.All() {
		if server.AppAddr != `` {
			for _, ifcAddr := range ifcAddrs {
				if strings.HasPrefix(ifcAddr.String(), server.Addr+`/`) {
					return server
				}
			}
		}
	}
	return Server{}
}

func (s *serverVar) CurrentTasks() string {
	tasks := []string{}
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range s.All() {
	loop:
		for _, ifcAddr := range ifcAddrs {
			if strings.HasPrefix(ifcAddr.String(), server.Addr+`/`) {
				tasks = append(tasks, server.Tasks)
				break loop
			}
		}
	}
	return strings.Join(tasks, ` `)
}

func (s serverVar) Matched(feature string) []Server {
	matched := []Server{}
	for _, server := range s.All() {
		if strings.Contains(server.Tasks, feature) ||
			strings.Contains(server.Addr, feature) {
			matched = append(matched, server)
		}
	}
	return matched
}

func (s serverVar) MatchedAddrs(feature string) []string {
	addrs := []string{}
	for _, server := range s.All() {
		if strings.Contains(server.Tasks, feature) ||
			strings.Contains(server.Addr, feature) {
			if !contains(addrs, server.Addr) {
				addrs = append(addrs, server.Addr)
			}
		}
	}
	return addrs
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
