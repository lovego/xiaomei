package web

import (
	"path/filepath"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PortEnvName() string {
	return `NGPORT`
}

func (i Image) Envs() []string {
	if conf.Type() != `simple` {
		return nil
	}
	ports := simpleconf.PortsOf(`app`)
	if len(ports) == 0 {
		return nil
	}
	addrs := []string{}
	for _, port := range ports {
		addrs = append(addrs, `127.0.0.1:`+port)
	}
	return []string{`NGBackendAddrs=` + strings.Join(addrs, `,`)}
}

func (i Image) FilesForRun() []string {
	root := filepath.Join(release.Root(), `img-web`)
	return []string{
		root + `/site.conf.tmpl:/etc/nginx/sites-available/` + release.Name() + `.conf.tmpl`,
		root + `/public:/var/www/` + release.Name(),
	}
}
