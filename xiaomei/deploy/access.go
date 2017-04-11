package deploy

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func AccessPrint() error {
	accessConf, err := getAccessConf()
	if err != nil {
		return err
	}
	print(accessConf)
	return nil
}

func AccessSetup() error {
	script := fmt.Sprintf(`
	sudo tee /etc/nginx/sites-enabled/%s.conf > /dev/null &&
	sudo mkdir -p /var/log/nginx/%s &&
	sudo nginx -t &&
	sudo service nginx restart
	`, release.Name(), release.Name(),
	)
	accessConf, err := getAccessConf()
	if err != nil {
		return err
	}
	for _, node := range cluster.AccessNodes() {
		if _, err := node.Run(
			cmd.O{Stdin: strings.NewReader(accessConf)}, script,
		); err != nil {
			return err
		}
	}
	return nil
}

func getAccessConf() (string, error) {
	tmpl := template.Must(template.New(``).Parse(accessConfTmpl))
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
	svcName := getServiceToAccess()
	addrs := getDriver().AccessAddrs(svcName)
	if len(addrs) == 1 {
		data.BackendAddr = addrs[0]
	} else {
		data.UpstreamName = release.Name() + `_` + svcName
		data.UpstreamAddrs = addrs
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

func getServiceToAccess() string {
	services := conf.ServiceNames()
	if services[`web`] {
		return `web`
	}
	if services[`app`] {
		return `app`
	}
	panic(`no backend service found in ` + conf.File() + `.`)
}

const accessConfTmpl = `
{{- if .UpstreamName -}}
upstream {{ .UpstreamName }} {
  {{- range .UpstreamAddrs }}
  server {{ . }};
  {{- end }}
}
{{ end -}}

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
