package deploy

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

// TODO: keep container history
const deployScriptTmpl = `set -e
{{ range .VolumesToCreate }}
docker volume create {{ . }}
{{- end }}

deploy() {
  local name=$1
  local args=$2
  docker stop $name >/dev/null 2>&1 && docker rm $name
  id=$(docker run --name=$name -d --network=host --restart=always $args)
  while status=$(docker ps -f id="$id" --format {{ "'{{.Status}}'" }}); do
    echo "$name: $status"
    case "$status" in
      "Up "*" (health: starting)" ) sleep 1 ;;
      "Up "*                      ) break   ;;
           *                      ) docker logs "$id"; exit  1 ;;
    esac
  done
}

{{ range .Services -}}
args='{{range .Envs}}-e {{.}} {{end}}{{.Options}} {{.Image}} {{.Command}}'
{{ $svc := . -}}
{{ range .Instances -}}
deploy {{$svc.Name}}.{{.}} "-e {{$svc.InstanceEnvName}}={{.}} $args"
{{ else -}}
deploy {{.Name}} "$args"
{{ end }}
{{ end -}}
`

func getDeployScript(svcNames []string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConfig(svcNames)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

type deployConfig struct {
	VolumesToCreate []string
	Services        []serviceConfig
}
type serviceConfig struct {
	Name, Image, InstanceEnvName, Command, Options string
	Instances, Envs                                []string
}

func getDeployConfig(svcNames []string) deployConfig {
	data := deployConfig{
		VolumesToCreate: conf.Get().VolumesToCreate,
	}
	for _, svcName := range svcNames {
		data.Services = append(data.Services, getServiceConf(svcName))
	}
	return data
}

func getServiceConf(svcName string) serviceConfig {
	image := images.Get(svcName)
	service := conf.GetService(svcName)
	data := serviceConfig{
		Name:            release.DeployName() + `_` + svcName,
		Image:           image.NameWithDigestInRegistry(),
		InstanceEnvName: image.InstanceEnvName(),
		Envs:            image.Envs(),
		Command:         strings.Join(service.Command, ` `),
		Options:         strings.Join(service.Options, ` `),
	}
	if data.InstanceEnvName != `` {
		data.Instances = conf.InstancesOf(svcName)
	}
	return data
}
