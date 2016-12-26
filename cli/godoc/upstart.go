package godoc

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Upstart(gopath, godoc string) error {
	if err := writeUpstartConfig(gopath, godoc); err != nil {
		return err
	}
	cmd.Run(cmd.O{}, `sudo`, `stop`, `godoc`)
	_, err := cmd.Run(cmd.O{}, `sudo`, `start`, `godoc`)
	return err
}

func writeUpstartConfig(gopath, godoc string) error {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ GoPath, Godoc, AddrPort, IndexInterval string }{
		GoPath:        gopath,
		Godoc:         strings.TrimSpace(godoc),
		AddrPort:      config.Servers.CurrentAppServer().GodocAddr(),
		IndexInterval: config.Godoc.IndexInterval(),
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/godoc.conf`, &buf)
	return nil
}
