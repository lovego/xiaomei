package host

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/release"
)

const deployScriptTmpl = `
set -e
deploy() {
	name={{.Name}}{{ if .Ports -}} .$1 {{- end }}
	docker stop $name >/dev/null 2>&1 && docker rm $name
  docker run --name=$name {{ if .Ports }}-e={{.PortEnv}}=$1{{ end }} \
	{{ range .Envs }} -e {{ . }}{{ end }} \
	{{ range .Volumes}} -v {{ . }}{{ end }} \
	--network=host --restart=always {{.Image}}
}
docker volume create {{ .LogsVolume }}
{{ if .Ports }} for port in {{ .Ports }}; do deploy $port; done {{ else }} deploy {{ end }}
`

// TODO: keep container history
func getDeployScript(svcName string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConfig(svcName)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

type deployConf struct {
	Name, Image, Ports, PortEnv, LogsVolume string
	Envs, Volumes                           []string
}

func getDeployConfig(svcName string) deployConf {
	logsVolume := release.Name() + `_logs`
	conf := deployConf{
		Name:       release.Name() + `_` + svcName,
		Image:      Driver.ImageNameOf(svcName),
		Ports:      strings.Join(Driver.PortsOf(svcName), ` `),
		LogsVolume: logsVolume,
		Volumes:    getDeployVolumes(svcName, logsVolume),
	}
	switch svcName {
	case `app`:
		conf.PortEnv = `GOPORT`
		conf.Envs = []string{`GOENV=` + release.Env()}
	case `web`, `access`:
		conf.PortEnv = `NGPORT`
	default:
	}
	return conf
}

func getDeployVolumes(svcName, logsVolume string) []string {
	var volumes []string
	switch svcName {
	case `app`:
		volumes = []string{fmt.Sprintf(`%s:/home/ubuntu/%s/log`, logsVolume, release.Name())}
	case `web`, `access`:
		volumes = []string{fmt.Sprintf(`%s:/var/log/nginx`, logsVolume)}
	default:
	}
	return append(volumes, getService(svcName).Volumes...)
}
