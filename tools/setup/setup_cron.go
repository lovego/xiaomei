package setup

import (
	"bytes"
	"fmt"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"os/user"
	"path"
	"text/template"
)

type cronConfData struct {
	UserName, AppRoot, Env string
}

func SetupCron() {
	tmpl := template.Must(template.ParseFiles(
		path.Join(config.Root, `config/conf/cron.tmpl`),
	))

	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cronConfData{
		curUser.Username, config.Root, config.Data.Env,
	}); err != nil {
		panic(err)
	}

	deployName := config.Data.AppName + `_` + config.Data.Env

	cmd.SudoWriteFile(`/etc/cron.d/`+deployName, &buf)
	fmt.Println(`setup cron ok.`)
}
