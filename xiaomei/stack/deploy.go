package stack

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"text/template"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/images"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/fatih/color"
)

func Deploy(svcName string, doBuild, doPush bool) error {
	if doBuild {
		if err := images.Build(svcName); err != nil {
			return err
		}
	}
	if doPush {
		if err := images.Push(svcName); err != nil {
			return err
		}
	}
	if svcName == `` {
		utils.Log(color.GreenString(`deploying all services.`))
	} else {
		utils.Log(color.GreenString(`deploying ` + svcName + ` service.`))
	}
	stack, err := getDeployStack(svcName)
	if err != nil {
		return err
	}
	script, err := getDeployScript(svcName)
	if err != nil {
		return err
	}
	return cluster.Run(cmd.O{Stdin: bytes.NewReader(stack)}, script)
}

func getDeployStack(svcName string) ([]byte, error) {
	stack := release.GetStack()
	if svcName != `` {
		stack.Services = map[string]release.Service{svcName: stack.Services[svcName]}
	}
	for svcName, service := range stack.Services {
		if svcName == `app` {
			service[`environment`] = map[string]string{`GOENV`: release.Env()}
		}
	}
	return yaml.Marshal(stack)
}

const deployScriptTmpl = `
	cd && mkdir -p {{ .DirName }} && cd {{ .DirName }} &&
	cat - > {{ .FileName }}.yml &&
	docker stack deploy --compose-file={{ .FileName }}.yml {{ .Name }}
`

func getDeployScript(svcName string) (string, error) {
	deployConf := struct {
		Name, DirName, FileName string
	}{
		Name: release.Name(), DirName: release.Name() + `_` + release.Env(), FileName: `stack`,
	}
	if svcName != `` {
		deployConf.FileName = svcName
	}

	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, deployConf); err != nil {
		return ``, err
	}
	return buf.String(), nil
}
