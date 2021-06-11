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
	// set -e
	//
	// docker volume create example-logs >/dev/null
	// test $(uname) = Linux && isLinux=true || isLinux=false
	//
	// deploy() {
	//   local name=$1
	//   local args=$2
	//   local portEnvVar=$3
	//   local port=$4
	//
	//   $isLinux && args=" --network=host $args"
	//   if test -n "$portEnvVar"; then
	//     args="-e $portEnvVar=$port $args"
	//     $isLinux || args="-p $port:$port $args"
	//   fi
	//
	//   docker stop $name >/dev/null 2>&1 && docker rm $name >/dev/null
	//   id=$(docker run --name=$name -dt --restart=always $args)
	//   echo $name
	//   docker logs -f $id 2>&1 | { sed '/started./q'; pkill -P $$ docker; }
	// }
	//
	// args='-e ProENV=production -v example-logs:/home/ubuntu/logs registry.example.com/example/app:production-180803-141210'
	// deploy app.3001 "$args" "ProPORT" 3001
	// deploy app.4001 "$args" "ProPORT" 4001
	// args='-e ProENV=production -v example-logs:/home/ubuntu/example-logs registry.example.com/example/logc:production-180803-141210'
	// deploy logc "$args"
}
