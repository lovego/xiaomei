package deploy

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
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
	isRollback := commit != ``
	tag, err := setupDeployTag(commit)
	if err != nil {
		return err
	}

	servers := config.Servers.Matched(serverFilter)

	var updated = make(map[string]bool)
	for _, server := range servers {
		sshAddr := server.SshAddr()
		if !updated[server.Addr] {
			updateCodeAndBin(sshAddr, tag, isRollback)
		}
		setupServer(sshAddr, tag, server.Tasks)
		updated[server.Addr] = true
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
	return nil
}

var updateTmpl *template.Template

func updateCodeAndBin(sshAddr, tag string, isRollback bool) {
	gitAddr := config.Deploy.GitAddr()
	if gitAddr == `` {
		panic(`no such git address`)
	}
	gitHost := getGitHost(gitAddr)
	updateConf := UpdateConfig{
		AppName:    config.App.Name(),
		DeployPath: config.Deploy.Path(),
		GitBranch:  config.Deploy.GitBranch(),
		GitTag:     tag,
		GitHost:    gitHost,
		GitAddr:    gitAddr,
	}
	color.Cyan(sshAddr)

	if updateTmpl == nil {
		updateTmpl = template.Must(template.New(``).Parse(updateCodeShell))
	}

	var buf bytes.Buffer
	err := updateTmpl.Execute(&buf, updateConf)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, sshAddr, buf.String())
	if !isRollback {
		cmd.Run(cmd.O{Panic: true}, `scp`, path.Join(config.App.Root(), updateConf.AppName),
			sshAddr+`:`+path.Join(updateConf.DeployPath, `release/bins`, tag))
	}
}

var deployTmpl *template.Template

func setupServer(sshAddr, tag string, tasks []string) {
	deployConf := DeployConfig{
		AppName:    config.App.Name(),
		DeployPath: config.Deploy.Path(),
		Env:        config.App.Env(),
		Tasks:      strings.Join(tasks, ` `),
		GitTag:     tag,
	}
	color.Cyan(sshAddr)

	if deployTmpl == nil {
		deployTmpl = template.Must(template.New(``).Parse(setupShell))
	}

	var buf bytes.Buffer
	err := deployTmpl.Execute(&buf, deployConf)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, sshAddr, buf.String())
}

func getGitHost(gitAddr string) string {
	re := regexp.MustCompile(`@(.*):`)
	return re.FindStringSubmatch(gitAddr)[1]
}
