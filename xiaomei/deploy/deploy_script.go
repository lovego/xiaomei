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
  id=$(docker run --name=$name -d --restart=always $args)
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
args='{{.CommonArgs}}'
{{ $svc := . -}}
{{ range .Instances -}}
deploy {{$svc.Name}}.{{.}} "-e {{$svc.InstanceEnvName}}={{.}} $args"
{{ else -}}
deploy {{.Name}} "$args"
{{ end }}
{{ end -}}
`

func getDeployScript(env string, svcNames []string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConfig(env, svcNames)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

type deployConfig struct {
	VolumesToCreate []string
	Services        []serviceConfig
}
type serviceConfig struct {
	Name, InstanceEnvName, CommonArgs string
	Instances                         []string
}

func getDeployConfig(env string, svcNames []string) deployConfig {
	data := deployConfig{
		VolumesToCreate: conf.Get(env).VolumesToCreate,
	}
	for _, svcName := range svcNames {
		data.Services = append(data.Services, getServiceConf(env, svcName))
	}
	return data
}

func getServiceConf(env, svcName string) serviceConfig {
	commonArgs := getCommonArgs(env, svcName, true)
	data := serviceConfig{
		Name:            release.ServiceName(env, svcName),
		InstanceEnvName: images.Get(svcName).InstanceEnvName(),
		CommonArgs:      strings.Join(commonArgs, ` `),
	}
	if data.InstanceEnvName != `` {
		data.Instances = conf.GetService(env, svcName).Instances()
	}
	return data
}
