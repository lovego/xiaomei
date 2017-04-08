package access

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func Config() (string, error) {
	tmpl := template.Must(template.New(``).Parse(configTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getConfigData()); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

func getConfigData() interface{} {
	data := struct {
		ProName, ServerNames, BackendAddr, UpstreamName string
		UpstreamAddrs                                   []string
	}{
		ProName:     release.Name(),
		ServerNames: getServerNames(),
	}
	svcName := getBackendSvcName()
	backendName := release.Name() + `_` + svcName
	if port := publicPorts(svcName); port == `` {
		data.BackendAddr = backendName + `:` + portOfService(svcName)
	} else {
		data.UpstreamName = backendName
		data.UpstreamAddrs = getBackendAddrs(port)
	}
	return data
}

func getServerNames() string {
	result := []string{}
	for _, env := range cluster.Envs() {
		result = append(result, release.AppIn(env).Domain())
	}
	return strings.Join(result, ` `)
}

func getBackendSvcName() string {
	services := deploy.ServiceNames()
	if services[`web`] {
		return `web`
	}
	if services[`app`] {
		return `app`
	}
	panic(`no backend service found in release.yml.`)
}

func publicPorts(svcName string) string {
	if ports := deploy.PortsOf(svcName); len(ports) > 0 {
		port := ports[0]
		return port[:strings.IndexByte(port, ':')]
	}
	return ``
}

func getBackendAddrs(port string) (addrs []string) {
	for _, node := range cluster.AccessNodes() {
		addrs = append(addrs, node.GetListenAddr()+`:`+port)
	}
	return
}

func portOfService(svcName string) string {
	switch svcName {
	case `app`:
		return `3000`
	case `web`:
		return `80`
	default:
		panic(`unexpected svcName: ` + svcName)
	}
}

const configTmpl = `
{{ if .UpstreamName -}}
upstream {{ .UpstreamName }} {
  {{- range .UpstreamAddrs }}
  server {{ . }};
  {{- end }}
}
{{- end }}

server {
  listen 80;
  server_name {{ .ServerNames }};

  location / {
    proxy_pass   http://{{ or .BackendAddr .UpstreamName }};
    include proxy_params;
  }

  access_log /var/log/nginx/{{ .ProName }}/access.log std;
  error_log  /var/log/nginx/{{ .ProName }}/access.err;
}
`
