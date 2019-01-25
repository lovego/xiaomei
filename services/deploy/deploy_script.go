package deploy

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/services/deploy/conf"
	"github.com/lovego/xiaomei/services/images"
	"github.com/lovego/xiaomei/release"
)

const deployScriptTmpl = `set -e
{{ range .VolumesToCreate }}
docker volume create {{ . }} >/dev/null
{{- end }}
test $(uname) = Linux && isLinux=true || isLinux=false

deploy() {
  local name=$1
  local args=$2
  local portEnvVar=$3
  local port=$4

  $isLinux && args=" --network=host $args"
  if test -n "$portEnvVar"; then
    args="-e $portEnvVar=$port $args"
    $isLinux || args="-p $port:$port $args"
  fi

  docker stop $name >/dev/null 2>&1 && docker rm $name >/dev/null
  id=$(docker run --name=$name -dt --restart=always $args)
  echo $name
  docker logs -f $id 2>&1 | { sed '/started./q'; pkill -P $$ docker; }
}

{{ range .Services -}}
  args='{{ .CommonArgs }}'{{"\n"}}
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
		VolumesToCreate: conf.Get(env).VolumesToCreate,
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
		data.Ports = conf.GetService(svcName, env).Ports
	}
	return data
}
