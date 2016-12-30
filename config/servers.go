package config

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
	for _, server := range s.All() {
		if server.HasTask(`appserver`) && server.IsLocal() {
			return &server
		}
	}
	return nil
}

func (s *ServerConf) CurrentTasks() (tasks []string) {
	for _, server := range s.All() {
		if server.IsLocal() {
			tasks = append(tasks, server.Tasks...)
		}
	}
	return
}

func (s *ServerConf) Matched(feature string) []Server {
	if feature == `` {
		return s.All()
	}
	matched := []Server{}
	for _, server := range s.All() {
		if server.Match(feature) {
			matched = append(matched, server)
		}
	}
	return matched
}

func (s *ServerConf) Matched2(feature, task string) []Server {
	if feature == `` && task == `` {
		return s.All()
	}
	matched := []Server{}
	for _, server := range s.All() {
		if server.Match(feature) && server.HasTask(task) {
			matched = append(matched, server)
		}
	}
	return matched
}
