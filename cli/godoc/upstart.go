package godoc

import (
	"bytes"
	"text/template"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

type upstartConf struct {
	GoPath, GodocBin, Addr, IndexInterval string
}

func setupUpstart(conf *upstartConf) error {
	if err := writeUpstartConfig(conf); err != nil {
		return err
	}
	cmd.Run(cmd.O{}, `sudo`, `stop`, `godoc`)
	_, err := cmd.Run(cmd.O{}, `sudo`, `start`, `godoc`)
	return err
}

func writeUpstartConfig(conf *upstartConf) error {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conf); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/godoc.conf`, &buf)
	return nil
}
