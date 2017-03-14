package access

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Config() (string, error) {
	var buf bytes.Buffer
	if err := template.Must(template.New(``).Parse(configTmpl)).Execute(&buf, struct {
		ServerNames, BackendAddr, ProName string
	}{
		ServerNames: getServerNames(),
		BackendAddr: getBackendAddr(),
		ProName:     release.Name(),
	}); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

func getServerNames() string {
	result := []string{}
	for _, env := range release.Envs() {
		result = append(result, release.AppIn(env).Domain())
	}
	return strings.Join(result, ` `)
}

func getBackendAddr() string {
	svcName := accessSvcName()
	if ports := release.PortsOf(svcName); len(ports) > 0 {
		port := ports[0]
		port = port[0:strings.IndexByte(port, ':')]
		return `127.0.0.1:` + port
	}
	switch svcName {
	case `app`:
		return release.Name() + `_app:3000`
	case `web`:
		return release.Name() + `_web:80`
	default:
		panic(`unexpected svcName: ` + svcName)
	}
}

func accessSvcName() string {
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
  server_name {{ .ServerNames }};

  location / {
    proxy_pass   http://{{ .BackendAddr }};
    include proxy_params;
  }

  access_log /var/log/nginx/{{ .ProName }}/access.log std;
  error_log  /var/log/nginx/{{ .ProName }}/access.err;
}
`
