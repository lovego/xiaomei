package web

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PortEnvName() string {
	return `NGINXPORT`
}

func (i Image) EnvsForRun() []string {
	return []string{`SendfileOff=true`}
}

func (i Image) FilesForRun() []string {
	root := filepath.Join(release.Root(), `img-web`)
	return []string{
		root + `/web.conf.tmpl:/etc/nginx/sites-available/` + release.Name() + `.conf.tmpl`,
		root + `/public:/var/www/` + release.Name(),
	}
}
