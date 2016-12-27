package appserver

import (
	"bytes"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Setup() {
	writeUpstartConfig()
	Restart()
}

func writeUpstartConfig() {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config.Data()); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/`+config.Deploy.Name()+`.conf`, &buf)
}
