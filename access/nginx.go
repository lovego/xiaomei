package access

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/release"
)

var reloadScript = `
sudo nginx -t
sudo nginx -s reload
`

var setupScriptTmpl = template.Must(template.New(``).Parse(`
set -eu
sudo tee /etc/nginx/sites-enabled/{{ .AccessName }} > /dev/null
sudo mkdir -p /var/log/nginx/{{ .AccessName }}
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
	return clusterRun(env, feature, "", "set -eu\n"+reloadScript)
}

func SetupNginx(env, feature, downAddr string) error {
	nginxConf, data, err := getNginxConf(env, downAddr)
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

	output, err := misc.RenderFile(filepath.Join(release.Root(env), `access.conf.tmpl`), data)
	if err != nil {
		return ``, Config{}, err
	}

	return output.String(), data, nil
}

func clusterRun(env, feature, input, cmdStr string) error {
	accessNodes := release.GetDeploy(env).AccessNodes
	for _, node := range release.GetCluster(env).GetNodes(feature) {
		if node.Match(accessNodes) {
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

func nodes(env string) error {
	accessNodes := release.GetDeploy(env).AccessNodes
	for _, node := range release.GetCluster(env).GetNodes("") {
		if node.Match(accessNodes) {
			fmt.Println(node.SshAddr())
		}
	}
	return nil
}
