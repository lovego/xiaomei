package setup

import (
	"bytes"
	"fmt"
	"os/user"
	"path"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"text/template"
)

type cronConfData struct {
	UserName, AppRoot, Env string
}

func SetupCron() {
	tmpl := template.Must(template.ParseFiles(
		path.Join(config.Root, `config/conf/cron.tmpl`),
	))

	current_user, err := user.Current()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cronConfData{
		current_user.Username, config.Root, config.Data.Env,
	}); err != nil {
		panic(err)
	}

	deploy_name := config.Data.AppName + `_` + config.Data.Env

	cmd.SudoWriteFile(`/etc/cron.d/`+deploy_name, &buf)
	fmt.Println(`setup cron ok.`)
}
