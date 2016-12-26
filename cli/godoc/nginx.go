package godoc

import (
	"bytes"
	"path"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Nginx() {
	writeNginxConfig()

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
}

func writeNginxConfig() {
	tmpl := template.Must(template.New(``).Parse(nginxConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getConfData()); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(path.Join(`/etc/nginx/sites-enabled/`, config.Deploy.Name()), &buf)
}

type confData struct {
	DeployName, AppRoot, AppPort, Domain string
	Servers                              []config.Server
	Nfs                                  bool
}

func getConfData() confData {
	fs, _ := cmd.Run(cmd.O{Panic: true, Output: true},
		`stat`, `--file-system`, `--format`, `%T`, config.App.Root(),
	)
	return confData{
		DeployName: config.Deploy.Name(),
		AppRoot:    config.App.Root(),
		AppPort:    config.App.Port(),
		Domain:     config.App.Domain(),
		Servers:    config.Servers.All(),
		Nfs:        fs == `nfs`,
	}
}
