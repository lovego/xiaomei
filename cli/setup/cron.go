package setup

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

type cronConfData struct {
	UserName, AppRoot, Env string
}

func SetupCron() {
	filePath := path.Join(config.App.Root(), `deploy/cron.tmpl`)
	if !utils.IsFile(filePath) {
		fmt.Println(`no such file: ` + filePath)
		return
	}
	tmpl := template.Must(template.ParseFiles(filePath))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cronConfData{
		config.Deploy.User(), config.App.Root(), config.App.Env(),
	}); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/cron.d/`+config.Deploy.Name(), &buf)
	fmt.Println(`setup cron ok.`)
}
