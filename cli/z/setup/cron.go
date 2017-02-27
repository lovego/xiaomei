package setup

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/fs"
)

func SetupCron() {
	filePath := path.Join(config.App.Root(), `deploy/cron.tmpl`)
	if !fs.IsFile(filePath) {
		fmt.Println(`no such file: ` + filePath)
		return
	}
	tmpl := template.Must(template.ParseFiles(filePath))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config.Data()); err != nil {
		panic(err)
	}

	cmd.SudoWriteFile(`/etc/cron.d/`+config.Deploy.Name(), &buf)
	fmt.Println(`setup cron ok.`)
}
