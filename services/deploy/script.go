package deploy

import (
	"bytes"
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
    case $(dockerStatus $name) in
      '' )
        ;;
      running )
        dockerRemove $name.old
        docker rename $name $name.old
        ;;
      * )
        dockerRemove $name
        ;;
    esac
    checkPort $port $name.old
  else
    dockerRemove $name
  fi
  docker run --name=$name -dt --restart=always $args
	docker logs -f $name |& { timeout ${StartTimeout:-1m} sed '/ started\./q'; pkill -P $$ docker; }

  test -n "$portEnvVar" && [[ $(dockerStatus $name.old) != '' ]] && dockerStop $name.old || true
}

dockerRemove() {
  [[ $(dockerStatus $1) == '' ]] && return
  dockerStop $1
  docker rm  $1
}

dockerStop() {
  time docker stop -t 180 $1 >/dev/null
}

dockerStatus() {
  docker inspect -f '{{ "{{ .State.Status }}" }}' $1 2>/dev/null || true
}

checkPort() {
  local port=$1

  local pidList=$(sudo lsof -itcp:$port -stcp:listen -Fp | grep -oP '^p\K\d+$')
  for pid in $pidList; do
    local dockerId=$(cat /proc/$pid/cgroup | grep -oP -m1 ':/docker/\K\w+$')
    if test -n "$dockerId"; then
      local container=$(docker inspect -f '{{ "{{ .Name }}" }}' $dockerId)
      container=${container#/}
      if [[ $container != $2 ]]; then
		echo "$port is already bound by other container: $container."
        sudo lsof -itcp:$port -stcp:listen -P
        exit 1
	  fi
    else
      echo "$port is already bound by other process."
      sudo lsof -itcp:$port -stcp:listen -P
      exit 1
    fi
  done
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
		Name:       release.ServiceName(env, svcName),
		CommonArgs: strings.Join(commonArgs, ` `),
		PortEnvVar: images.Get(svcName).PortEnvVar(),
	}
	if data.PortEnvVar != `` {
		data.Ports = release.GetService(env, svcName).Ports
	}
	return data
}
