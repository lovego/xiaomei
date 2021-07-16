package deploy

import (
	"bytes"
	"fmt"
	"text/template"
)

func testScriptTmpl(conf deployConfig) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	fmt.Println(tmpl.Execute(&buf, conf))
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
	// <nil>
}
