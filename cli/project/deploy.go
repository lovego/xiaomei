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
			var env string
			if len(args) > 1 {
				return errors.New(`redundant args.`)
			} else if len(args) == 1 {
				env = args[0]
			}
			return Deploy(env, ``)
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
	stack, err := GetStack()
	if err != nil {
		return nil, err
	}
	if svc != `` {
		stack.Services = map[string]Service{svc: stack.Services[svc]}
	}
	for _, service := range stack.Services {
		delete(service, `build`)
	}
	return yaml.Marshal(stack)
}

const deployScriptTmpl = `
	cd && mkdir -p {{ .DeployName }} && cd {{ .DeployName }} &&
	cat - > {{ .FileName }}.yml &&
	docker stack deploy --compose-file={{ .FileName }}.yml {{ .DeployName }}
`
