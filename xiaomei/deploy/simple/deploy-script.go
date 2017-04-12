package simple

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

// TODO: keep container history, wait until healthy
const deployScriptTmpl = `set -e
deploy() {
  name={{.Name}}{{ if .Ports }}.$1{{ end }}
  docker stop $name >/dev/null 2>&1 && docker rm $name
  docker run --name=$name {{ if .Ports }}-e {{.PortEnv}}=$1{{ end }} \
	{{ range .Envs }} -e {{ . }}{{ end }} \
  {{ range .Volumes}} -v {{ . }}{{ end }} \
  -d --network=host --restart=always \
	{{.Image}}
}
{{ range .VolumesToCreate -}}
docker volume create {{ . }}
{{ end }}
{{ if .Ports -}}
for port in {{ .Ports }}; do deploy $port; done
{{ else -}}
deploy
{{ end }}
exit
`

func getDeployScript(svcName string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConfig(svcName)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

type deployConf struct {
	Name, Image, PortEnv, Ports    string
	Envs, VolumesToCreate, Volumes []string
}

func getDeployConfig(svcName string) deployConf {
	image := images.Get(svcName)
	conf := deployConf{
		Name:            release.Name() + `_` + svcName,
		Image:           image.NameWithDigestInRegistry(),
		PortEnv:         image.PortEnvName(),
		Envs:            image.Envs(),
		VolumesToCreate: simpleconf.Get().VolumesToCreate,
		Volumes:         simpleconf.GetService(svcName).Volumes,
	}
	if conf.PortEnv != `` {
		conf.Ports = strings.Join(simpleconf.PortsOf(svcName), ` `)
	}
	return conf
}
