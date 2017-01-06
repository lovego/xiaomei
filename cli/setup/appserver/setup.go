package appserver

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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
	if err := tmpl.Execute(&buf, struct {
		config.Conf
		Path string
	}{
		config.Data(), upstartPath(),
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/apps/`+config.Deploy.Name()+`.conf`, &buf)
}

func upstartPath() string {
	path := os.Getenv(`PATH`)
	if path == `` {
		path = `/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin`
	}
	if gopath := os.Getenv(`GOPATH`); gopath != `` {
		for _, workspace := range strings.Split(gopath, `:`) {
			path += `:` + filepath.Join(workspace, `bin`)
		}
	}
	return path
}
