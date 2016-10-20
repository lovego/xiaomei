package setup

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
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

	current_user, err := user.Current()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, upstartConfData{
		config.Data, current_user.Username, config.Root, config.CurrentAppServer().AppStartOn,
	}); err != nil {
		panic(err)
	}

	deploy_name := config.Data.AppName + `_` + config.Data.Env

	cmd.SudoWriteFile(`/etc/init/`+deploy_name+`.conf`, &buf)

	cmd.Run(cmd.O{}, `sudo`, `stop`, deploy_name)
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `sudo`, `start`, deploy_name)
	fmt.Println(output)
	if !strings.Contains(output, `start/running,`) {
		os.Exit(1)
	}
}
