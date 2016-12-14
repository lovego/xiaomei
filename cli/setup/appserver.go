package setup

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func SetupAppServer() {
	writeUpstartConfig()

	// stop current
	cmd.Run(cmd.O{}, `sudo`, `stop`, config.DeployName())

	errLog := path.Join(config.Root(), `log/appserver.stderr`)
	cmd.Run(cmd.O{Panic: true}, `touch`, `-a`, errLog)
	tail, _ := cmd.Start(cmd.O{Panic: true}, `tail`, `-n0`, `-f`, errLog)
	// start new
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `sudo`, `start`, config.DeployName())
	tail.Process.Kill()

	fmt.Println(output)
	if !strings.Contains(output, `start/running,`) {
		os.Exit(1)
	}
}

type upstartConfData struct {
	UserName, AppRoot, AppPort, AppStartOn string
}

func writeUpstartConfig() {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, upstartConfData{
		UserName:   config.DeployUser(),
		AppRoot:    config.Root(),
		AppPort:    config.AppPort(),
		AppStartOn: config.CurrentAppServer().AppStartOn,
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/`+config.DeployName()+`.conf`, &buf)
}
