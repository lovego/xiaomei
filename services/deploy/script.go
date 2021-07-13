package deploy

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
)

const deployScriptTmpl = `set -e
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
    if [[ $(docker inspect -f '{{ "{{ .State.Status }}" }}' $name 2>/dev/null) == running ]]; then
      dockerRemove $name.old
      docker rename $name $name.old
    fi
    checkPort $port $name.old
  else
    dockerRemove $name
  fi
  set -x
  docker run --name=$name -dt --restart=always $args
  set +x
  docker logs -f $name |& { sed '/ started\./q'; pkill -P $$ docker; }

  test -n "$portEnvVar" && dockerStop $name.old
}

dockerStop() {
  if [[ $(docker inspect -f '{{ "{{ .State.Status }}" }}' $1 2>/dev/null) != running ]]; then
    return
  fi
  set -x
  time docker stop -t 180 $1 >/dev/null
  set +x
}

dockerRemove() {
  dockerStop $1
  docker rm  $1 &>/dev/null || true
}

checkPort() {
  local port=$1

  local pid=$(lsof -itcp:$port -stcp:listen -Fp | grep -oP '^p\K\d+$')
  test -z "$pid" && return
  local dockerId=$(cat /proc/$pid/cgroup | grep -oP -m1 ':/docker/\K\w+$')
  if test -n "$dockerId"; then
    local container=$(docker inspect -f '{{ "{{ .Name }}" }}' $dockerId)
    container=${container#/}
    [[ $container == $2 ]] && return
    echo "$port is already bound by container $container: "
  else
    echo "$port is already bound by: "
  fi
  lsof -itcp:$port -stcp:listen -P
  exit 1
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
	script := buf.String()
	// fmt.Println(script)
	return script, nil
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
