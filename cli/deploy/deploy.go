package deploy

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

type DeployConfig struct {
	DeployPath, Env, Tasks              string
	GitBranch, GitTag, GitHost, GitAddr string
}

func Deploy(commit, serverFilter string) error {
	if err := os.Chdir(config.App.Root()); err != nil {
		return err
	}
	tag, err := setupDeployTag(commit)
	if err != nil {
		return err
	}

	gitHost := getGitHost(config.Deploy.GitAddr())
	servers := config.Servers.Matched(serverFilter)
	for _, server := range servers {
		deployToServer(config.Deploy.User()+`@`+server.Addr, DeployConfig{
			DeployPath: config.Deploy.Path(),
			Env:        config.App.Env(),
			Tasks:      server.Tasks,
			GitBranch:  config.Deploy.GitBranch(),
			GitTag:     tag,
			GitHost:    gitHost,
			GitAddr:    config.Deploy.GitAddr(),
		})
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
	return nil
}

var deployTmpl *template.Template

func deployToServer(userAddr string, deployConf DeployConfig) {
	color.Cyan(userAddr)

	if deployTmpl == nil {
		deployTmpl = template.Must(template.New(``).Parse(deployShell))
	}

	var buf bytes.Buffer
	err := deployTmpl.Execute(&buf, deployConf)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, userAddr, buf.String())
}

func getGitHost(gitAddr string) string {
	re := regexp.MustCompile(`@(.*):`)
	return re.FindStringSubmatch(gitAddr)[1]
}
