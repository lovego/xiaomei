package nginx

import (
	"bytes"
	// "fmt"
	"path"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Setup() {
	writeMainConfig()
	writeServerConfig()

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
}

func writeMainConfig() {
	confFile := path.Join(config.App.Root(), `deploy/nginx.conf`)
	if utils.IsFile(confFile) {
		cmd.Run(cmd.O{Panic: true}, `sudo`, `cp`, confFile, `/etc/nginx/`)
	} else {
		cmd.SudoWriteFile(`/etc/nginx/nginx.conf`, strings.NewReader(defaultMainConfig))
	}
}

func writeServerConfig() {
	var tmpl *template.Template
	confFile := path.Join(config.App.Root(), `deploy/nginx.tmpl.conf`)
	if utils.IsFile(confFile) {
		tmpl = template.Must(template.ParseFiles(confFile))
	} else {
		tmpl = template.Must(template.New(``).Parse(defaultServerConfig))
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
