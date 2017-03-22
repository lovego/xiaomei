package godoc

import (
	"bytes"
	"text/template"

	"github.com/lovego/xiaomei/utils/cmd"
)

type nginxConf struct {
	Addrs  []string
	Domain string
}

func setupNginx(conf *nginxConf) error {
	if err := writeNginxConfig(conf); err != nil {
		return err
	}

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
	return nil
}

func writeNginxConfig(conf *nginxConf) error {
	tmpl := template.Must(template.New(``).Parse(nginxConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conf); err != nil {
		return err
	}

	cmd.SudoWriteFile(`/etc/nginx/sites-enabled/godoc`, &buf)
	return nil
}
