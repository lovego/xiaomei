package access

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

var nginxSetupScriptTmpl = template.Must(template.New(``).Parse(`
  set -e
  {{ if .Https }}
	# sudo mkir -p /var/www/letsencrypt
	# sudo certbot certonly --webroot -w /var/www/letsencrypt -d {{ .Domain }}
	{{ end }}
	sudo tee /etc/nginx/sites-enabled/{{ .Domain }}.conf > /dev/null
	sudo mkdir -p /var/log/nginx/{{ .Domain }}
	sudo nginx -t
	sudo service nginx reload || which reload-nginx && reload-nginx
`))

func setupNginx(env, feature, nginxConf string, data interface{}) error {
	var script bytes.Buffer
	if err := nginxSetupScriptTmpl.Execute(&script, data); err != nil {
		return err
	}

	for _, node := range cluster.Get(env).GetNodes(feature) {
		if node.Labels[`access`] == `true` {
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
