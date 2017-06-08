package deploy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/utils/fs"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func accessPrint(svcName string) error {
	accessConf, err := getAccessConf(svcName)
	if err != nil {
		return err
	}
	fmt.Print(accessConf)
	return nil
}

func accessSetup(svcName string) error {
	script := fmt.Sprintf(`
	sudo tee /etc/nginx/sites-enabled/%s.conf > /dev/null &&
	sudo mkdir -p /var/log/nginx/%s &&
	sudo nginx -t &&
	sudo service nginx reload
	`, getServiceName(svcName), release.DeployName(),
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
	var confTmpl string
	if path := filepath.Join(release.Root(), `access.conf.tmpl`); fs.IsFile(path) {
		if buf, err := ioutil.ReadFile(path); err == nil {
			confTmpl = string(buf)
		} else {
			return ``, err
		}
	} else {
		confTmpl = defaultAccessConfTmpl
	}
	tmpl := template.Must(template.New(``).Parse(confTmpl))
	configData, err := getConfigData(svcName)
	if err != nil {
		return ``, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, configData); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

func getConfigData(svcName string) (interface{}, error) {
	data := struct {
		SvcName, ServiceName, DomainNames, BackendAddr, UpstreamName string
		UpstreamAddrs                                                []string
	}{
		SvcName:     svcName,
		ServiceName: getServiceName(svcName),
		DomainNames: getDomainNames(svcName),
	}
	if svcName == `` {
		svcName = getServiceToAccess()
	}
	addrs, err := getDriver().Addrs(svcName)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 1 {
		data.BackendAddr = addrs[0]
	} else {
		data.UpstreamName = getServiceName(svcName)
		data.UpstreamAddrs = addrs
	}
	return data, nil
}

func getDomainNames(svcName string) string {
	result := []string{}
	for _, env := range cluster.Envs() {
		domain := release.AppIn(env).Domain()
		if svcName != `` {
			domain = getSvcDomain(domain, svcName)
		}
		result = append(result, domain)
	}
	return strings.Join(result, ` `)
}

func getSvcDomain(domain, svcName string) string {
	parts := strings.SplitN(domain, `.`, 2)
	if len(parts) == 2 {
		return parts[0] + `-` + svcName + `.` + parts[1]
	} else {
		return domain + `-` + svcName
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

func getServiceName(svcName string) string {
	name := release.DeployName()
	if svcName != `` {
		name += `_` + svcName
	}
	return name
}

const defaultAccessConfTmpl = `
{{- if .UpstreamName -}}
upstream {{ .UpstreamName }} {
  {{- range .UpstreamAddrs }}
  server {{ . }};
  {{- end }}
}
{{ end -}}

server {
  listen 80;
  server_name {{ .DomainNames }};

  location / {
    proxy_pass   http://{{ or .BackendAddr .UpstreamName }};
    include proxy_params;
  }

  access_log /var/log/nginx/{{ .ServiceName }}/access.log std;
  error_log  /var/log/nginx/{{ .ServiceName }}/access.err;
}
`
