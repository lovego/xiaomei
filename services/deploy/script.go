package deploy

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
)

const deployScriptTmpl = `set -ex
{{ range .VolumesToCreate }}
docker volume create {{ . }} >/dev/null
{{- end }}
if [[ $(uname) == Linux ]]; then
  isLinux=true
  networkArgs="--network=host"
else
  isLinux=false
fi

deploy() {
  local name=$1
  local args=$2
  local portEnvVar=$3
  local port=$4

  if test -n "$portEnvVar"; then
    args="-e $portEnvVar=$port $args"
    $isLinux || args="-p $port:$port $args"
    if [[ $(docker inspect -f '{{ "{{ .State.Status }}" }}' $name) == running ]]; then
      dockerRemove $name.old
      docker rename $name $name.old
    fi
    checkPort $port $name.old
  else
    dockerRemove $name
  fi
  docker run --name=$name -dt --restart=always $args
  docker logs -f $name |& { sed '/ started\./q'; pkill -P $$ docker; }

  test -n "$portEnvVar" && dockerRemove $name.old
}

dockerRemove() {
  docker stop $1 &>/dev/null || true
  docker rm   $1 &>/dev/null || true
}

checkPort() {
  true
}

{{ range .Services -}}
  args="$networkArgs "'{{ .CommonArgs }}'{{"\n"}}
  {{- $svc := . -}}
  {{ if .PortEnvVar -}}
    {{ range .Ports -}}
      deploy {{$svc.Name}}.{{.}} "$args" "{{$svc.PortEnvVar}}" {{.}}{{"\n"}}
    {{- end -}}
  {{- else -}}
    deploy {{$svc.Name}} "$args"{{"\n"}}
  {{- end -}}
{{ end -}}
`

func getDeployScript(svcNames []string, env, timeTag string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConfig(svcNames, env, timeTag)); err != nil {
		return ``, err
	}
	fmt.Println(buf.String())
	return buf.String(), nil
}

type deployConfig struct {
	VolumesToCreate []string
	Services        []serviceConfig
}
type serviceConfig struct {
	Name, CommonArgs, PortEnvVar string
	Ports                        []uint16
}

func getDeployConfig(svcNames []string, env, timeTag string) deployConfig {
	data := deployConfig{
		VolumesToCreate: release.GetDeploy(env).VolumesToCreate,
	}
	for _, svcName := range svcNames {
		data.Services = append(data.Services, getServiceConf(svcName, env, timeTag))
	}
	return data
}

func getServiceConf(svcName, env, timeTag string) serviceConfig {
	commonArgs := GetCommonArgs(svcName, env, timeTag)
	data := serviceConfig{
		Name:       release.ServiceName(svcName, env),
		CommonArgs: strings.Join(commonArgs, ` `),
		PortEnvVar: images.Get(svcName).PortEnvVar(),
	}
	if data.PortEnvVar != `` {
		data.Ports = release.GetService(env, svcName).Ports
	}
	return data
}
