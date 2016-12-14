package setup

import (
	"bytes"
	// "fmt"
	"path"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func SetupNginx() {
	writeNginxConfig()

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
}

func writeNginxConfig() {
	var tmpl *template.Template
	confFile := path.Join(config.Root(), `deploy/nginx.tmpl.conf`)
	if utils.IsFile(confFile) {
		tmpl = template.Must(template.ParseFiles(confFile))
	} else {
		tmpl = template.Must(template.New(``).Parse(nginxConfig))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getNginxConfData()); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(path.Join(`/etc/nginx/sites-enabled/`, config.DeployName()), &buf)
}

type nginxConfData struct {
	DeployName, AppRoot, AppPort, Domain string
	Nfs                                  bool
}

func getNginxConfData() nginxConfData {
	fs, _ := cmd.Run(cmd.O{Panic: true, Output: true},
		`stat`, `--file-system`, `--format`, `%T`, config.Root(),
	)
	return nginxConfData{
		DeployName: config.DeployName(),
		AppRoot:    config.Root(),
		AppPort:    config.AppPort(),
		Domain:     config.Domain(),
		Nfs:        fs == `nfs`,
	}
}
