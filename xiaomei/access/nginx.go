package access

import (
	"bytes"
	"fmt"
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

var reloadScript = `
if test -f /lib/systemd/system/nginx.service; then
	sudo systemctl reload nginx
else
	sudo service nginx reload
fi`

var setupScriptTmpl = template.Must(template.New(``).Parse(`
set -e
sudo tee /etc/nginx/sites-enabled/{{ .Domain }} > /dev/null
sudo mkdir -p /var/log/nginx/{{ .Domain }}
sudo nginx -t
` + reloadScript))

func HasAccess(svcs []string) bool {
	for _, svcName := range svcs {
		if svcName == "app" || svcName == "web" || svcName == "godoc" {
			return true
		}
	}
	return false
}

func ReloadNginx(env, feature string) error {
	return clusterRun(env, feature, "", reloadScript)
}

func SetupNginx(env, svcName, feature, downAddr string) error {
	data, err := getConfig(env, svcName, downAddr)
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
	return clusterRun(env, feature, nginxConf, script.String())
}

func clusterRun(env, feature, input, cmdStr string) error {
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if node.Labels[`access`] == `true` {
			log.Println(color.GreenString(node.SshAddr()))
			cmdOpt := cmd.O{}
			if input != "" {
				cmdOpt.Stdin = strings.NewReader(input)
			}
			if _, err := node.Run(cmdOpt, cmdStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func printNginxConf(env, svcName string) error {
	data, err := getConfig(env, svcName, "")
	if err != nil {
		return err
	}
	nginxConf, err := getNginxConf(svcName, data)
	if err != nil {
		return err
	}
	fmt.Print(nginxConf)
	return nil
}

func getNginxConf(svcName string, data interface{}) (string, error) {
	name := svcName
	if name == `` {
		name = `access`
	}
	file := filepath.Join(release.Root(), name+`.conf.tmpl`)
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
