package deploy

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"text/template"

	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

type DeployConfig struct {
	Tasks, Addr     string
	GitTag, GitHost string
}

func Deploy(commit string) error {
	if err := os.Chdir(config.Root()); err != nil {
		return err
	}
	tag, err := setupDeployTag(commit)
	if err != nil {
		return err
	}

	gitHost := getGitHost(config.Data().GitAddr)
	servers := cli.MatchedServers()
	for _, server := range servers {
		deployToServer(DeployConfig{
			Tasks:   server.Tasks,
			Addr:    config.DeployUser + `@` + server.Addr,
			GitTag:  tag,
			GitHost: gitHost,
		})
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
	return nil
}

var deployTmpl *template.Template

func deployToServer(server DeployConfig) {
	color.Cyan(server.Addr)

	if deployTmpl == nil {
		deployTmpl = template.Must(template.ParseFiles(
			path.Join(config.Root(), `config/shell/deploy.tmpl.sh`),
		))
	}

	var buf bytes.Buffer
	err := deployTmpl.Execute(&buf, server)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, server.Addr, buf.String())
}

func getGitHost(gitAddr string) string {
	re := regexp.MustCompile(`@(.*):`)
	return re.FindStringSubmatch(gitAddr)[1]
}
