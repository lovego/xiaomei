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

func accessPrint(svcName string) error {
	accessConf, err := getAccessConf(svcName)
	if err != nil {
		return err
	}
	print(accessConf)
	return nil
}

func accessSetup(svcName string) error {
	name := release.Name()
	if svcName != `` {
		name += `_` + svcName
	}
	script := fmt.Sprintf(`
	sudo tee /etc/nginx/sites-enabled/%s.conf > /dev/null &&
	sudo mkdir -p /var/log/nginx/%s &&
	sudo nginx -t &&
	sudo service nginx reload
	`, name, release.Name(),
	)
	accessConf, err := getAccessConf(svcName)
	if err != nil {
		return err
	}
	for _, node := range cluster.Nodes() {
		if node.Labels[`access`] == `true` {
			if _, err := node.Run(
				cmd.O{Stdin: strings.NewReader(accessConf)}, script,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func getAccessConf(svcName string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(accessConfTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getConfigData(svcName)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

func getConfigData(svcName string) interface{} {
	data := struct {
		ProName, ServerNames, BackendAddr, UpstreamName string
		UpstreamAddrs                                   []string
	}{
		ProName:     release.Name(),
		ServerNames: getServerNames(svcName),
	}
	if svcName == `` {
		svcName = getServiceToAccess()
	}
	addrs := getDriver().AccessAddrs(svcName)
	if len(addrs) == 1 {
		data.BackendAddr = addrs[0]
	} else {
		data.UpstreamName = release.Name() + `_` + svcName
		data.UpstreamAddrs = addrs
	}
	return data
}

func getServerNames(svcName string) string {
	result := []string{}
	for _, env := range cluster.Envs() {
		domain := release.AppIn(env).Domain()
		if svcName == `godoc` {
			domain = godocDomain(domain)
		}
		result = append(result, domain)
	}
	return strings.Join(result, ` `)
}

func godocDomain(domain string) string {
	parts := strings.SplitN(domain, `.`, 2)
	if len(parts) == 2 {
		return parts[0] + `-godoc.` + parts[1]
	} else {
		return domain + `-godoc`
	}
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
