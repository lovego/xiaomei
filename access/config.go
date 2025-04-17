package access

import (
	"errors"
	"strconv"

	"github.com/lovego/config/config"
	"github.com/lovego/xiaomei/release"
)

type Config struct {
	AccessName string
	*config.Config
	App, Web *service
}

func getConfig(env, downAddr string) (Config, error) {
	data := Config{
		Config: release.Config(env),
		App:    newService(`app`, env, downAddr),
		Web:    newService(`web`, env, downAddr),
	}
	if data.App == nil && data.Web == nil {
		return Config{}, errors.New(`neither app nor web service defined.`)
	}
	if data.App != nil && data.Web != nil {
		data.AccessName = data.Config.DeployName()
	} else if data.App != nil {
		data.AccessName = release.ServiceName(env, data.App.svcName)
	} else if data.Web != nil {
		data.AccessName = release.ServiceName(env, data.Web.svcName)
	}
	return data, nil
}

type service struct {
	*config.Config
	svcName  string
	downAddr string
}

func newService(svcName, env, downAddr string) *service {
	if release.HasService(env, svcName) {
		return &service{Config: release.Config(env), svcName: svcName, downAddr: downAddr}
	} else {
		return nil
	}
}

func (s *service) Addrs(feature ...string) ([]string, error) {
	if s == nil {
		return nil, nil
	}
	addrs := []string{}
	ports := release.GetService(s.Env.String(), s.svcName).Ports
	for _, node := range s.Nodes(feature...) {
		for _, port := range ports {
			upstreamAddr := node.GetServiceAddr() + `:` + strconv.FormatInt(int64(port), 10)
			if s.downAddr != "" && s.downAddr == node.Addr {
				upstreamAddr += " down"
			}
			addrs = append(addrs, upstreamAddr)
		}
	}
	if len(addrs) == 0 {
		return nil, errors.New(`no instance defined for: ` + s.svcName)
	}
	return addrs, nil
}

func (s *service) Nodes(feature ...string) (nodes []release.Node) {
	if s == nil {
		return nil
	}
	var f string
	if len(feature) > 0 {
		f = feature[0]
	}
	labels := release.GetService(s.Env.String(), s.svcName).Nodes
	for _, node := range release.GetCluster(s.Env.String()).GetNodes(f) {
		if node.Match(labels) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
