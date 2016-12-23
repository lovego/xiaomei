package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

type UpdateConfig struct {
	AppName, DeployPath, GitBranch, GitTag, GitHost, GitAddr string
}

type DeployConfig struct {
	AppName, DeployPath, Env, Tasks, GitTag string
}

func Deploy(commit, serverFilter string) error {
	if err := os.Chdir(config.App.Root()); err != nil {
		return err
	}
	isRollback := false
	if commit != `` {
		isRollback = true
	}
	tag, err := setupDeployTag(commit)
	if err != nil {
		return err
	}

	gitAddr := config.Deploy.GitAddr()
	if gitAddr == `` {
		return errors.New(`no such git address`)
	}
	gitHost := getGitHost(gitAddr)
	servers := config.Servers.Matched(serverFilter)

	var updated = make(map[string]bool)
	for _, server := range servers {
		userAddr := config.Deploy.User() + `@` + server.Addr
		if !updated[server.Addr] {
			updateCodeAndBin(userAddr, isRollback, UpdateConfig{
				AppName:    config.App.Name(),
				DeployPath: config.Deploy.Path(),
				GitBranch:  config.Deploy.GitBranch(),
				GitTag:     tag,
				GitHost:    gitHost,
				GitAddr:    gitAddr,
			})
		}
		deployToServer(userAddr, DeployConfig{
			AppName:    config.App.Name(),
			DeployPath: config.Deploy.Path(),
			Env:        config.App.Env(),
			Tasks:      server.Tasks,
			GitTag:     tag,
		})
		updated[server.Addr] = true
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
	return nil
}

var updateTmpl *template.Template

func updateCodeAndBin(userAddr string, isRollback bool, updateConf UpdateConfig) {
	color.Cyan(userAddr)

	if updateTmpl == nil {
		updateTmpl = template.Must(template.New(``).Parse(updateShell))
	}

	var buf bytes.Buffer
	err := updateTmpl.Execute(&buf, updateConf)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, userAddr, buf.String())
	if !isRollback {
		cmd.Run(cmd.O{Panic: true}, `scp`, path.Join(config.App.Root(), updateConf.AppName),
			userAddr+`:`+path.Join(updateConf.DeployPath, `release/bins`, updateConf.GitTag))
	}
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
