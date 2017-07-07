package access

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/release"
)

func getNginxConf(svcName string) (string, string, error) {
	tmpl, err := getNginxConfTmpl(svcName)
	if err != nil {
		return ``, ``, err
	}
	data := getConfData(svcName)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return ``, ``, err
	}
	return buf.String(), data.Domain(), nil
}

func getNginxConfTmpl(name string) (*template.Template, error) {
	if name == `` {
		name = `access`
	}
	file := filepath.Join(release.Root(), `access`, name+`.conf.tmpl`)
	confTmpl, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return template.Must(template.New(``).Parse(string(confTmpl))), nil
}

type domain interface {
	Domain() string
}

func getConfData(svcName string) domain {
	if svcName == `` {
		return accessConfig{
			App: newService(`app`),
			Web: newService(`web`),
		}
	} else {
		return newService(svcName)
	}
}

type accessConfig struct {
	App, Web *service
}

func (a accessConfig) Env() string {
	return release.Env()
}

func (a accessConfig) Domain() string {
	return release.App().Domain()
}
