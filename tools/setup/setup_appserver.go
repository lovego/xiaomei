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
	tmpl := template.Must(template.ParseFiles(
		path.Join(config.Root, `config/conf/upstart.tmpl.conf`),
	))

	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, upstartConfData{
		config.Data, curUser.Username, config.Root, config.CurrentAppServer().AppStartOn,
	}); err != nil {
		panic(err)
	}

	deployName := config.Data.AppName + `_` + config.Data.Env

	cmd.SudoWriteFile(`/etc/init/`+deployName+`.conf`, &buf)

	cmd.Run(cmd.O{}, `sudo`, `stop`, deployName)
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `sudo`, `start`, deployName)
	fmt.Println(output)
	if !strings.Contains(output, `start/running,`) {
		os.Exit(1)
	}
}
