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
	"github.com/lovego/xiaomei/release"
)

var reloadScript = `
sudo nginx -t
sudo systemctl reload nginx
`

var setupScriptTmpl = template.Must(template.New(``).Parse(`
set -e
sudo tee /etc/nginx/sites-enabled/{{ .Domain }} > /dev/null
sudo mkdir -p /var/log/nginx/{{ .Domain }}
` + reloadScript))

func HasAccess(svcs []string) bool {
	for _, svcName := range svcs {
		if svcName == "app" || svcName == "web" {
			return true
		}
	}
	return false
}

func ReloadNginx(env, feature string) error {
	return clusterRun(env, feature, "", reloadScript)
}

func SetupNginx(env, feature, downAddr string) error {
	nginxConf, data, err := getNginxConf(env, "")
	if err != nil {
		return err
	}
	var script bytes.Buffer
	if err := setupScriptTmpl.Execute(&script, data); err != nil {
		return err
	}
	return clusterRun(env, feature, nginxConf, script.String())
}

func printNginxConf(env string) error {
	nginxConf, _, err := getNginxConf(env, "")
	if err != nil {
		return err
	}
	fmt.Print(nginxConf)
	return nil
}

func getNginxConf(env, downAddr string) (string, Config, error) {
	data, err := getConfig(env, downAddr)
	if err != nil {
		return ``, Config{}, err
	}

	file := filepath.Join(release.Root(), `access.conf.tmpl`)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return ``, Config{}, err
	}
	tmpl := template.Must(template.New(``).Parse(string(content)))
	if err != nil {
		return ``, Config{}, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return ``, Config{}, err
	}
	return buf.String(), data, nil
}

func clusterRun(env, feature, input, cmdStr string) error {
	access := release.GetService(`access`, env)
	for _, node := range release.GetCluster(env).GetNodes(feature) {
		if node.Match(access.Nodes) {
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
