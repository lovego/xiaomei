package config

import (
	"net"
	"strings"
)

var Servers ServerConf

type ServerConf struct {
	conf serversConf
}

type serversConf []Server

func (s *ServerConf) All() []Server {
	Load()
	return s.conf
}

func (s *ServerConf) CurrentAppServer() *Server {
	addrs := machineAddrs()
	for _, server := range s.All() {
		if server.HasTask(`appserver`) && server.HasAddr(addrs) {
			return &server
		}
	}
	return nil
}

func (s *ServerConf) CurrentTasks() (tasks []string) {
	addrs := machineAddrs()
	for _, server := range s.All() {
		if server.HasAddr(addrs) {
			tasks = append(tasks, server.Tasks...)
		}
	}
	return
}

func (s *ServerConf) Matched(feature string) []Server {
	matched := []Server{}
	for _, server := range s.All() {
		if server.Match(feature) {
			matched = append(matched, server)
		}
	}
	return matched
}

func (s *ServerConf) MatchedAppserver(feature string) []Server {
	matched := []Server{}
	for _, server := range s.All() {
		if server.Match(feature) && server.HasTask(`appserver`) {
			matched = append(matched, server)
		}
	}
	return matched
}

func machineAddrs() []string {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	result := make([]string, len(ifcAddrs))
	for i, ifcAddr := range ifcAddrs {
		addr := ifcAddr.String()
		if i := strings.IndexByte(addr, '/'); i >= 0 {
			addr = addr[:i]
		}
		result[i] = addr
	}
	return result
}
