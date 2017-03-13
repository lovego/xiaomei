package access

import (
	"bytes"
	"text/template"

	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func PrintConfig(svcName string) error {
	svcName = accessSvcName(svcName)
	data := struct {
		ProName, SvcName, BackendPort string
	}{ProName: release.Name(), SvcName: release.Name() + `_` + svcName}
	switch svcName {
	case `app`:
		data.BackendPort = `3000`
	case `web`:
		data.BackendPort = `80`
	}
	tmpl := template.Must(template.New(``).Parse(configTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	println(buf.String())
	return nil
}

func accessSvcName(svcName string) string {
	if svcName != `` {
		return svcName
	}
	stack := release.GetStack()
	if stack.Services[`web`] != nil {
		return `web`
	}
	if stack.Services[`app`] != nil {
		return `app`
	}
	panic(`no app service found in stack.yml.`)
}

const configTmpl = `
server {
  listen 80;

  location / {
    proxy_pass   http://{{ .SvcName }}:{{ .BackendPort }};
    include proxy_params;
  }

  access_log /var/log/nginx/{{ .ProName }}/access.log std;
  error_log  /var/log/nginx/{{ .ProName }}/access.err;
}
`
