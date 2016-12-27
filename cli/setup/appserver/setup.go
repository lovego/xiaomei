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

type upstartConfData struct {
	UserName, AppRoot, AppName, AppStartOn string
}

func writeUpstartConfig() {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, upstartConfData{
		UserName:   config.Deploy.User(),
		AppRoot:    config.App.Root(),
		AppName:    config.App.Name(),
		AppStartOn: config.Servers.CurrentAppServer().AppStartOn,
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/`+config.Deploy.Name()+`.conf`, &buf)
}
