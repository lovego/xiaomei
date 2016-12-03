package setup

import (
	"bytes"
	"fmt"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"os"
	"os/user"
	"path"
	"strings"
	"text/template"
)

type upstartConfData struct {
	*config.Config
	UserName, AppRoot, AppStartOn string
}

func SetupAppServer() {
	writeUpstartConfig()

	// stop current
	cmd.Run(cmd.O{}, `sudo`, `stop`, config.Data.DeployName)

	errLog := path.Join(config.Root(), `log/app.err`)
	cmd.Run(cmd.O{Panic: true}, `touch`, `-a`, errLog)
	tail, _ := cmd.Start(cmd.O{Panic: true}, `tail`, `-n0`, `-f`, errLog)
	// start new
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `sudo`, `start`, config.Data.DeployName)
	tail.Process.Kill()

	fmt.Println(output)
	if !strings.Contains(output, `start/running,`) {
		os.Exit(1)
	}
}

func writeUpstartConfig() {
	tmpl := template.Must(template.ParseFiles(
		path.Join(config.Root(), `deploy/conf/upstart.tmpl.conf`),
	))

	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, upstartConfData{
		config.Data, curUser.Username, config.Root(), config.CurrentAppServer().AppStartOn,
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/init/`+config.Data.DeployName+`.conf`, &buf)
}
