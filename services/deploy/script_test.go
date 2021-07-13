package deploy

import (
	"bytes"
	"fmt"
	"text/template"
)

var testTmpl = template.Must(template.New(``).Parse(deployScriptTmpl))

func testScriptTmpl(conf deployConfig) {
	var buf bytes.Buffer
	if err := testTmpl.Execute(&buf, conf); err != nil {
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
	//     if [[ $(docker inspect -f '{{ .State.Status }}' $name) == running ]]; then
	//       dockerRemove $name.old
	//       docker rename $name $name.old
	//     fi
	//     checkPort $port $name.old
	//   else
	//     dockerRemove $name
	//   fi
	//   docker run --name=$name -dt --restart=always $args
	//   docker logs -f $name |& { sed '/ started\./q'; pkill -P $$ docker; }
	//
	//   test -n "$portEnvVar" && dockerRemove $name.old
	// }
	//
	// dockerRemove() {
	//   docker stop $1 &>/dev/null || true
	//   docker rm   $1 &>/dev/null || true
	// }
	//
	// checkPort() {
	//   true
	// }
	//
	// args="$networkArgs "'-e ProENV=production -v example-logs:/home/ubuntu/logs registry.example.com/example/app:production-180803-141210'
	// deploy app.3001 "$args" "ProPORT" 3001
	// deploy app.4001 "$args" "ProPORT" 4001
	// args="$networkArgs "'-e ProENV=production -v example-logs:/home/ubuntu/example-logs registry.example.com/example/logc:production-180803-141210'
	// deploy logc "$args"
}
