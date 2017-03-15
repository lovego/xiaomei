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

func Deploy(svcName string, noBuild, noPush bool) error {
	if !noBuild {
		if err := images.Build(svcName); err != nil {
			return err
		}
	}
	if !noPush {
		if err := images.Push(svcName); err != nil {
			return err
		}
	}
	if err := deploy(svcName); err != nil {
		return err
	}
	return Ps(svcName, nil, true)
}

func deploy(svcName string) error {
	if svcName == `` {
		utils.Log(color.GreenString(`deploying all services.`))
	} else {
		utils.Log(color.GreenString(`deploying ` + svcName + ` service.`))
	}
	stackYaml, err := getDeployStack(svcName)
	if err != nil {
		return err
	}
	script, err := getDeployScript(svcName)
	if err != nil {
		return err
	}
	return cluster.Run(cmd.O{Stdin: bytes.NewReader(stackYaml)}, script)
}

func getDeployStack(svcName string) ([]byte, error) {
	stack := release.GetStack()
	if svcName != `` {
		stack.Services = map[string]release.Service{svcName: release.GetService(svcName)}
	}
	if app, ok := stack.Services[`app`]; ok {
		app[`environment`] = map[string]string{`GOENV`: release.Env()}
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
		Name: release.Name(), DirName: release.Name() + `_` + release.Env(), FileName: svcName,
	}

	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, deployConf); err != nil {
		return ``, err
	}
	return buf.String(), nil
}
