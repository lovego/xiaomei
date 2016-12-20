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
	cmd.Run(cmd.O{}, `sudo`, `stop`, config.Deploy.Name())

	appserverLog := path.Join(config.App.Root(), `log/appserver.log`)
	cmd.Run(cmd.O{Panic: true}, `touch`, `-a`, appserverLog)
	tail, _ := cmd.Start(cmd.O{Panic: true}, `tail`, `-n0`, `-f`, appserverLog)
	// start new
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `sudo`, `start`, config.Deploy.Name())
	tail.Process.Kill()

	fmt.Println(output)
	if !strings.Contains(output, `start/running,`) {
		os.Exit(1)
	}
}

type upstartConfData struct {
	UserName, AppRoot, AppName, AppPort, AppStartOn string
	StartTimeout                                    uint16
}

func writeUpstartConfig() {
	tmpl := template.Must(template.New(``).Parse(upstartConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, upstartConfData{
		UserName:     config.Deploy.User(),
		AppRoot:      config.App.Root(),
		AppName:      config.App.Name(),
		AppPort:      config.App.Port(),
		AppStartOn:   config.Servers.Current().AppStartOn,
		StartTimeout: config.App.StartTimeout(),
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/`+config.Deploy.Name()+`.conf`, &buf)
}
