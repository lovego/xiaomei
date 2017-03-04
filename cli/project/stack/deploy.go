package stack

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"text/template"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Deploy(svcName string) error {
	if err := Build(svcName); err != nil {
		return err
	}
	if err := push(svcName); err != nil {
		return err
	}
	if svcName == `` {
		config.Log(color.GreenString(`deploying all services.`))
	} else {
		config.Log(color.GreenString(`deploying ` + svcName + ` service.`))
	}
	stack, err := getDeployStack(svcName)
	if err != nil {
		return err
	}
	script, err := getDeployScript(svcName)
	if err != nil {
		return err
	}
	return cluster.Run(cmd.O{Stdin: bytes.NewReader(stack)}, config.Env(), script)
}

func getDeployStack(svcName string) ([]byte, error) {
	stack, err := getStack()
	if err != nil {
		return nil, err
	}
	if svcName != `` {
		stack.Services = map[string]Service{svcName: stack.Services[svcName]}
	}
	for svcName, service := range stack.Services {
		if imageName, err := ImageName(svcName); err != nil {
			return nil, err
		} else {
			service[`image`] = imageName
		}
		if svcName == `app` {
			service[`environment`] = map[string]string{`GOENV`: config.Env()}
			service[`volumes`] = []string{`./log:/home/ubuntu/appserver/log`}
		}
	}
	return yaml.Marshal(stack)
}

const deployScriptTmpl = `
	cd && mkdir -p {{ .DeployName }}/log && cd {{ .DeployName }} &&
	cat - > {{ .FileName }}.yml &&
	docker stack deploy --compose-file={{ .FileName }}.yml {{ .Name }}
`

func getDeployScript(svcName string) (string, error) {
	deployConf := struct {
		Name, DeployName, FileName string
	}{
		Name: config.Name(), DeployName: config.DeployName(), FileName: `stack`,
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

func push(svcName string) error {
	if svcName == `` {
		return eachServiceDo(push)
	}
	config.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	imageName, err := ImageName(svcName)
	if err != nil {
		return err
	}
	_, err = cmd.Run(cmd.O{}, `docker`, `push`, imageName)
	return err
}
