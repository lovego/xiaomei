package access

import (
	"bytes"
	"fmt"
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
	data, err := getConfData(svcName)
	if err != nil {
		return ``, ``, err
	}

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

func getConfData(svcName string) (domain, error) {
	if svcName == `` {
		data := accessConfig{
			App: newService(`app`),
			Web: newService(`web`),
		}
		if data.App == nil && data.Web == nil {
			return nil, fmt.Errorf(`neither app nor web service defined.`, svcName)
		}
		return data, nil
	} else {
		data := newService(svcName)
		if data == nil {
			return nil, fmt.Errorf(`%s service not defined.`, svcName)
		}
		return data, nil
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
