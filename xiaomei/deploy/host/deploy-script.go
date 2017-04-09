package host

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

const deployScriptTmpl = `
set -e
deploy() {
	name={{.Name}}.$1
	docker stop $name >/dev/null 2>&1 && docker rm $name
  docker run --name=$name -e {{.PortEnv}}=$1 \
	{{ range .Envs }} -e {{ . }}{{ end }} \
	{{ range .Volumes}} -v {{ . }}{{ end }} \
	--network=host --restart=always -d {{.Image}}
}
{{ range .VolumesToCreate }}
docker volume create {{ .Volume }}
{{ end }}
for port in {{ .Ports }}; do deploy $port; done
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
	Name, Image, Ports, PortEnv    string
	Envs, VolumesToCreate, Volumes []string
}

func getDeployConfig(svcName string) deployConf {
	return deployConf{
		Name:            release.Name() + `_` + svcName,
		Image:           Driver.ImageNameOf(svcName),
		Ports:           strings.Join(portsOf(svcName), ` `),
		PortEnv:         portEnvName(svcName),
		Envs:            images.Get(svcName).EnvsForDeploy(),
		VolumesToCreate: getRelease().VolumesToCreate,
		Volumes:         getService(svcName).Volumes,
	}
}
