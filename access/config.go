package access

import (
	"errors"
	"strconv"

	"github.com/lovego/config/config"
	"github.com/lovego/strmap"
	"github.com/lovego/xiaomei/release"
)

type Config struct {
	AccessName string
	*config.EnvConfig
	App, Web *service
	Data     strmap.StrMap
}

func getConfig(env, downAddr string) (Config, error) {
	data := Config{
		EnvConfig: release.EnvConfig(env),
		App:       newService(`app`, env, downAddr),
		Web:       newService(`web`, env, downAddr),
		Data:      release.EnvData(env),
	}
	if data.App == nil && data.Web == nil {
		return Config{}, errors.New(`neither app nor web service defined.`)
	}
	if data.App != nil && data.Web != nil {
		data.AccessName = data.EnvConfig.DeployName()
	} else if data.App != nil {
		data.AccessName = release.ServiceName(data.App.svcName, env)
	} else if data.Web != nil {
		data.AccessName = release.ServiceName(data.Web.svcName, env)
	}
	return data, nil
}

type service struct {
	*config.EnvConfig
	svcName  string
	downAddr string
	addrs    []string
}

func newService(svcName, env, downAddr string) *service {
	if release.HasService(svcName, env) {
		return &service{EnvConfig: release.EnvConfig(env), svcName: svcName, downAddr: downAddr}
	} else {
		return nil
	}
}

func (s *service) Addrs() ([]string, error) {
	if s == nil {
		return nil, nil
	}
	if s.addrs == nil {
		addrs := []string{}
		ports := release.GetService(s.Env.String(), s.svcName).Ports
		for _, node := range s.Nodes() {
			for _, port := range ports {
				upstreamAddr := node.GetServiceAddr() + `:` + strconv.FormatInt(int64(port), 10)
				if s.downAddr != "" && s.downAddr == node.Addr {
					upstreamAddr += " down"
				}
				addrs = append(addrs, upstreamAddr)
			}
		}
		s.addrs = addrs
		if len(addrs) == 0 {
			return nil, errors.New(`no instance defined for: ` + s.svcName)
		}
	}
	return s.addrs, nil
}

func (s *service) Nodes() (nodes []release.Node) {
	if s == nil {
		return nil
	}
	labels := release.GetService(s.Env.String(), s.svcName).Nodes
	for _, node := range release.GetCluster(s.Env.String()).GetNodes(``) {
		if node.Match(labels) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
