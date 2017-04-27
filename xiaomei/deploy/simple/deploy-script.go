package simple

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf/simpleconf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

// TODO: keep container history
const deployScriptTmpl = `set -e
{{ range .VolumesToCreate }}
docker volume create {{ . }}
{{- end }}

deploy() {
  name=$1
  args=$2
  docker stop $name >/dev/null 2>&1 && docker rm $name
  id=$(docker run --name=$name -d --network=host --restart=always $args)
  while status=$(docker ps -f id="$id" --format {{ "'{{.Status}}'" }}); do
    echo "$name: $status"
    case "$status" in
      "Up "*" (health: starting)" ) sleep 1 ;;
      "Up "*                      ) break   ;;
           *                      ) exit  1 ;;
    esac
  done
}

{{ range .Services -}}
args='{{range .Envs}}-e {{.}} {{end}}{{range .Volumes}}-v {{.}}{{end}} {{.Image}} {{.Command}}'
{{ $svc := . -}}
{{ range .Ports -}}
deploy {{$svc.Name}}.{{.}} "-e {{$svc.PortEnv}}={{.}} $args"
{{ else -}}
deploy {{.Name}} "$args"
{{ end }}
{{ end -}}
`

func getDeployScript(svcNames []string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConf(svcNames)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

type deployConf struct {
	VolumesToCreate []string
	Services        []serviceConf
}
type serviceConf struct {
	Name, Image, PortEnv, Command string
	Ports, Envs, Volumes          []string
}

func getDeployConf(svcNames []string) deployConf {
	conf := deployConf{
		VolumesToCreate: simpleconf.Get().VolumesToCreate,
	}
	for _, svcName := range svcNames {
		conf.Services = append(conf.Services, getServiceConf(svcName))
	}
	return conf
}

func getServiceConf(svcName string) serviceConf {
	image := images.Get(svcName)
	service := simpleconf.GetService(svcName)
	conf := serviceConf{
		Name:    release.Name() + `_` + svcName,
		Image:   image.NameWithDigestInRegistry(),
		PortEnv: image.PortEnvName(),
		Envs:    image.Envs(),
		Command: strings.Join(service.Command, ` `),
		Volumes: service.Volumes,
	}
	if conf.PortEnv != `` {
		conf.Ports = simpleconf.PortsOf(svcName)
	}
	return conf
}
