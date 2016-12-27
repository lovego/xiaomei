package nginx

import (
	"bytes"
	// "fmt"
	"path"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Setup() {
	writeConfig()

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
}

func writeConfig() {
	var tmpl *template.Template
	confFile := path.Join(config.App.Root(), `deploy/nginx.tmpl.conf`)
	if utils.IsFile(confFile) {
		tmpl = template.Must(template.ParseFiles(confFile))
	} else {
		tmpl = template.Must(template.New(``).Parse(defaultConfig))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getConfData()); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(path.Join(`/etc/nginx/sites-enabled/`, config.Deploy.Name()), &buf)
}

type confData struct {
	config.Conf
	Nfs bool
}

func getConfData() confData {
	fs, _ := cmd.Run(cmd.O{Panic: true, Output: true},
		`stat`, `--file-system`, `--format`, `%T`, config.App.Root(),
	)
	return confData{
		Conf: config.Data(),
		Nfs:  fs == `nfs`,
	}
}
