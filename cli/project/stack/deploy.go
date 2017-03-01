package stack

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"text/template"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Deploy(env, svcName string) error {
	addr, err := cluster.ManagerAddr(env)
	if err != nil {
		return err
	}
	stack, err := getDeployStack(svcName)
	if err != nil {
		return err
	}
	script, err := getDeployScript(svcName)
	if err != nil {
		return err
	}
	_, err = cmd.SshRun(cmd.O{Stdin: bytes.NewReader(stack)}, addr, script)
	return err
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
		service[`image`] = stack.ImageName(svcName)
	}
	return yaml.Marshal(stack)
}

func getDeployScript(svcName string) (string, error) {
	deployConf := struct {
		DeployName, FileName string
	}{
		DeployName: config.DeployName(), FileName: `stack`,
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

const deployScriptTmpl = `
	cd && mkdir -p {{ .DeployName }} && cd {{ .DeployName }} &&
	cat - > {{ .FileName }}.yml &&
	docker stack deploy --compose-file={{ .FileName }}.yml {{ .DeployName }}
`
