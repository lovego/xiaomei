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
  id=$(
    docker run --name=$name -d --network=host --restart=always {{ if .Ports -}}
    -e {{.PortEnv}}=$1 {{ end }} {{ range .Envs -}} -e {{ . }} {{ end }} {{ range .Volumes -}}
    -v {{ . }} {{ end -}} {{.Image}} {{.Command}}
  )
  while status=$(docker ps -f id="$id" --format {{ "'{{.Status}}'" }}); do
    echo "$name: $status"
    case "$status" in
      "Up "*" (health: starting)" ) sleep 1 ;;
      "Up "*                      ) break ;;
           *                      ) exit 1;;
    esac
  done
}
{{ range .VolumesToCreate -}}
docker volume create {{ . }}
{{ end -}}
{{ if .Ports -}}
for port in {{ .Ports }}; do deploy $port; done
{{ else -}}
deploy
{{ end -}}
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
	Name, Image, PortEnv, Ports, Command string
	Envs, VolumesToCreate, Volumes       []string
}

func getDeployConfig(svcName string) deployConf {
	image := images.Get(svcName)
	service := simpleconf.GetService(svcName)
	conf := deployConf{
		Name:            release.Name() + `_` + svcName,
		Image:           image.NameWithDigestInRegistry(),
		PortEnv:         image.PortEnvName(),
		Envs:            image.Envs(),
		Command:         strings.Join(service.Command, ` `),
		VolumesToCreate: simpleconf.Get().VolumesToCreate,
		Volumes:         service.Volumes,
	}
	if conf.PortEnv != `` {
		conf.Ports = strings.Join(simpleconf.PortsOf(svcName), ` `)
	}
	return conf
}
