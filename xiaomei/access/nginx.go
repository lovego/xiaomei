package access

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

var setupScriptTmpl = template.Must(template.New(``).Parse(`
  set -e
	{{ if .CheckCert -}}
	test -e /letsencrypt/live/{{ .Domain }} && exit 0
	{{- end }}
	sudo tee /etc/nginx/sites-enabled/{{ .Domain }}.conf > /dev/null
	sudo mkdir -p /var/log/nginx/{{ .Domain }}
	sudo nginx -t
	sudo service nginx reload || { which reload-nginx && sudo reload-nginx; }
	{{ if .CheckCert -}}
	sudo mkdir -p /var/www/letsencrypt
	sudo certbot certonly --webroot -w /var/www/letsencrypt -d {{ .Domain }}
	{{- end }}
`))

func setupNginx(env, svcName, feature string, checkCert bool) error {
	data, err := getConfig(env, svcName, checkCert)
	if err != nil {
		return err
	}
	nginxConf, err := getNginxConf(svcName, data)
	if err != nil {
		return err
	}
	var script bytes.Buffer
	if err := setupScriptTmpl.Execute(&script, data); err != nil {
		return err
	}

	for _, node := range cluster.Get(env).GetNodes(feature) {
		if node.Labels[`access`] == `true` {
			log.Println(color.GreenString(node.SshAddr()))
			if _, err := node.Run(
				cmd.O{Stdin: strings.NewReader(nginxConf)}, script.String(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func getNginxConf(svcName string, data interface{}) (string, error) {
	name := svcName
	if name == `` {
		name = `access`
	}
	file := filepath.Join(release.Root(), `access`, name+`.conf.tmpl`)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return ``, err
	}
	tmpl := template.Must(template.New(``).Parse(string(content)))
	if err != nil {
		return ``, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return ``, err
	}
	return buf.String(), nil
}
