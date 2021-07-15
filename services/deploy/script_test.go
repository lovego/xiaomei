package deploy

import (
	"bytes"
	"fmt"
	"text/template"
)

func testScriptTmpl(conf deployConfig) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conf); err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(buf.String())
}

func ExampleDeployScript() {
	testScriptTmpl(deployConfig{
		VolumesToCreate: []string{"example-logs"},
		Services: []serviceConfig{
			{
				Name: "app",
				CommonArgs: "-e ProENV=production -v example-logs:/home/ubuntu/logs " +
					"registry.example.com/example/app:production-180803-141210",
				PortEnvVar: "ProPORT",
				Ports:      []uint16{3001, 4001},
			},
			{
				Name: "logc",
				CommonArgs: "-e ProENV=production -v example-logs:/home/ubuntu/example-logs " +
					"registry.example.com/example/logc:production-180803-141210",
			},
		},
	})
	// Output:
	// set -ex
	//
	// docker volume create example-logs >/dev/null
	// if [[ $(uname) == Linux ]]; then
	//   isLinux=true
	//   networkArgs="--network=host"
	// else
	//   isLinux=false
	// fi
	//
	// deploy() {
	//   local name=$1
	//   local args=$2
	//   local portEnvVar=$3
	//   local port=$4
	//
	//   if test -n "$portEnvVar"; then
	//     args="-e $portEnvVar=$port $args"
	//     $isLinux || args="-p $port:$port $args"
	//     case $(dockerStatus $name) in
	//       '' )
	//         ;;
	//       running )
	//         dockerRemove $name.old
	//         docker rename $name $name.old
	//         ;;
	//       * )
	//         dockerRemove $name
	//         ;;
	//     esac
	//     checkPort $port $name.old
	//   else
	//     dockerRemove $name
	//   fi
	//   docker run --name=$name -dt --restart=always $args
	//   docker logs -f $name |& { sed '/ started\./q'; pkill -P $$ docker; }
	//
	//   test -n "$portEnvVar" && [[ $(dockerStatus $name.old) != '' ]] && dockerStop $name.old
	// }
	//
	// dockerRemove() {
	//   [[ $(dockerStatus $1) == '' ]] && return
	//   dockerStop $1
	//   docker rm  $1
	// }
	//
	// dockerStop() {
	//   time docker stop -t 180 $1 >/dev/null
	// }
	//
	// dockerStatus() {
	//   docker inspect -f '{{ .State.Status }}' $1 2>/dev/null
	// }
	//
	// checkPort() {
	//   local port=$1
	//
	//   local pid=$(sudo lsof -itcp:$port -stcp:listen -Fp | grep -oP '^p\K\d+$')
	//   test -z "$pid" && return
	//   local dockerId=$(cat /proc/$pid/cgroup | grep -oP -m1 ':/docker/\K\w+$')
	//   if test -n "$dockerId"; then
	//     local container=$(docker inspect -f '{{ .Name }}' $dockerId)
	//     container=${container#/}
	//     [[ $container == $2 ]] && return
	//     echo "$port is already bound by container $container: "
	//   else
	//     echo "$port is already bound by: "
	//   fi
	//   sudo lsof -itcp:$port -stcp:listen -P
	//   exit 1
	// }
	//
	// args="$networkArgs "'-e ProENV=production -v example-logs:/home/ubuntu/logs registry.example.com/example/app:production-180803-141210'
	// deploy app.3001 "$args" "ProPORT" 3001
	// deploy app.4001 "$args" "ProPORT" 4001
	// args="$networkArgs "'-e ProENV=production -v example-logs:/home/ubuntu/example-logs registry.example.com/example/logc:production-180803-141210'
	// deploy logc "$args"
}
