package develop

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Godoc() error {
	if err := writeGodocUpstartConfig(); err != nil {
		return err
	}
	cmd.Run(cmd.O{}, `sudo`, `stop`, `godoc`)
	_, err := cmd.Run(cmd.O{}, `sudo`, `start`, `godoc`)
	return err
}

func writeGodocUpstartConfig() error {
	gopath := os.Getenv(`GOPATH`)
	if gopath == `` {
		return errors.New(`empty GOPATH.`)
	}
	godoc, _ := cmd.Run(cmd.O{Output: true, Panic: true}, `which`, `godoc`)

	tmpl := template.Must(template.New(``).Parse(`
start on runlevel [2345]
respawn
respawn limit 2 60

script
  GOPATH={{.GoPath}} {{.Godoc}} -http=:1234 -index_interval=1s
end script
`))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ GoPath, Godoc string }{
		GoPath: gopath,
		Godoc:  strings.TrimSpace(godoc),
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/godoc.conf`, &buf)
	return nil
}
