package project

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v2"
	"text/template"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func DeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `deploy <env>`,
		Short: `deploy project to the specified env.`,
		RunE: func(c *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return errors.New(`<env> is required.`)
			case 1:
				return Deploy(args[0], ``)
			default:
				return errors.New(`redundant args.`)
			}
		},
	}
}

func Deploy(env, svc string) error {
	addr, err := cluster.ManagerAddr(env)
	if err != nil {
		return err
	}
	stack, err := GetDeployStack(svc)
	if err != nil {
		return err
	}
	script, err := GetDeployScript(svc)
	if err != nil {
		return err
	}
	_, err = cmd.SshRun(cmd.O{Stdin: bytes.NewReader(stack)}, addr, script)
	return err
}

func GetDeployScript(svc string) (string, error) {
	deployConf := struct {
		DeployName, FileName string
	}{
		DeployName: config.DeployName(), FileName: `stack`,
	}
	if svc != `` {
		deployConf.FileName = svc
	}

	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, deployConf); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

func GetDeployStack(svc string) ([]byte, error) {
	if svc == `` {
		return GetStackFileContent()
	}
	stack, err := GetStack()
	if err != nil {
		return nil, err
	}
	stack.Services = map[string]Service{svc: stack.Services[svc]}
	return yaml.Marshal(stack)
}

const deployScriptTmpl = `
	cd && mkdir -p {{ .DeployName }} && cd {{ .DeployName }} &&
	cat - > {{ .FileName }}.yml &&
	docker stack deploy --compose-file={{ .FileName }}.yml {{ .DeployName }}
`
