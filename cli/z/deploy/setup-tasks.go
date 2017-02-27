package deploy

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func setupTasks(sshAddr, tag string, tasks []string) {
	script := getSetupTasksScript(tag, tasks)
	if config.IsLocalEnv() {
		cmd.Run(cmd.O{Panic: true}, `sh`, `-c`, script)
	} else {
		cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, sshAddr, script)
	}
}

var setupTmpl *template.Template

func getSetupTasksScript(tag string, tasks []string) string {
	deployConf := struct {
		config.Conf
		Tasks, GitTag string
	}{
		Conf:   config.Data(),
		Tasks:  strings.Join(tasks, ` `),
		GitTag: tag,
	}

	if setupTmpl == nil {
		setupTmpl = template.Must(template.New(``).Parse(setupTasksTmpl))
	}

	var buf bytes.Buffer
	err := setupTmpl.Execute(&buf, deployConf)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

const setupTasksTmpl = `
main () {
  cd {{.Deploy.Path}}/release || exit 1

  ln -sfT envs/{{.App.Env}}.yml config/env.yml || exit 1
  ln -sfT bins/{{ .GitTag }} {{ .App.Name }} || exit 1
  xiaomei setup {{ .Tasks }} || exit 1
  xiaomei deploy clear-local-tags || exit 1
}

main
`
