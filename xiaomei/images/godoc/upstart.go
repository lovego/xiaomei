package godoc

import (
	"bytes"
	"text/template"

	"github.com/lovego/xiaomei/utils/cmd"
)

type upstartConf struct {
	GoPath, GodocBin, Addr, IndexInterval string
}

func setupUpstart(conf *upstartConf) error {
	if err := writeUpstartConfig(conf); err != nil {
		return err
	}

	var buf bytes.Buffer
	cmd.Run(cmd.O{Stderr: &buf}, `sudo`, `stop`, `apps/godoc`)
	stdErr := buf.String()
	if stdErr != "stop: Unknown instance: \n" {
		print(stdErr)
	}

	_, err := cmd.Run(cmd.O{}, `sudo`, `start`, `apps/godoc`)
	return err
}

func writeUpstartConfig(conf *upstartConf) error {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conf); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/apps/godoc.conf`, &buf)
	return nil
}
