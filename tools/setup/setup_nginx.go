package setup

import (
	"bytes"
	// "fmt"
	"path"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

type nginxConfData struct {
	*config.Config
	AppRoot     string
	CurrentAddr string
	Nfs         bool
}

func SetupNginx() {
	data := getNginxConfData()

	writeNginxConfig(data)

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
}

func writeNginxConfig(data nginxConfData) {
	tmpl := template.Must(template.ParseFiles(
		path.Join(config.Root, `config/conf/nginx.tmpl.conf`),
	))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(path.Join(`/etc/nginx/sites-enabled/`, config.Data.DeployName), &buf)
}

func getNginxConfData() nginxConfData {
	fs, _ := cmd.Run(cmd.O{Panic: true, Output: true},
		`stat`, `--file-system`, `--format`, `%T`, config.Root)
	return nginxConfData{
		Config:      config.Data,
		AppRoot:     config.Root,
		CurrentAddr: config.CurrentAppServer().Addr,
		Nfs:         fs == `nfs`,
	}
}
