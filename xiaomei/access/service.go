package access

import (
	"errors"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

type service struct {
	svcName string
	addrs   []string
}

func newService(svcName string) *service {
	if conf.HasService(svcName) {
		return &service{svcName: svcName}
	} else {
		return nil
	}
}

func (s *service) Env() string {
	return release.Env()
}

func (s *service) Addrs() ([]string, error) {
	if s == nil {
		return nil, nil
	}
	if s.addrs == nil {
		addrs := []string{}
		instances := conf.InstancesOf(s.svcName)
		for _, node := range s.Nodes() {
			for _, instance := range instances {
				addrs = append(addrs, node.GetListenAddr()+`:`+instance)
			}
		}
		s.addrs = addrs
		if len(addrs) == 0 {
			return nil, errors.New(`no instance defined for: ` + s.svcName)
		}
	}
	return s.addrs, nil
}

func (s *service) Nodes() (nodes []cluster.Node) {
	if s == nil {
		return nil
	}
	service := conf.GetService(s.svcName)
	for _, node := range cluster.Nodes(``) {
		if node.Match(service.Nodes) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *service) Upstream() (string, error) {
	if s == nil {
		return ``, nil
	}
	if addrs, err := s.Addrs(); err != nil {
		return ``, err
	} else if len(addrs) > 1 {
		return release.DeployName() + `_` + s.svcName, nil
	} else {
		return ``, nil
	}
}

func (s *service) ProxyPass() (string, error) {
	if s == nil {
		return ``, nil
	}
	if addrs, err := s.Addrs(); err != nil {
		return ``, err
	} else if len(addrs) == 1 {
		return addrs[0], nil
	} else {
		return s.Upstream()
	}
}

func (s *service) Domain() string {
	if s == nil {
		return ``
	}
	domain := release.App().Domain()
	parts := strings.SplitN(domain, `.`, 2)
	if len(parts) == 2 {
		return parts[0] + `-` + s.svcName + `.` + parts[1]
	} else {
		return domain + `-` + s.svcName
	}
}
