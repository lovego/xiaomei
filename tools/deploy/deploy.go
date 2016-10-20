package deploy

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

type DeployConfig struct {
	*config.Config
	Tasks, Addr     string
	GitTag, GitHost string
}

func Deploy(tag string) {
	if err := os.Chdir(config.Root); err != nil {
		panic(err)
	}

	gitHost := getGitHost(config.Data.GitAddr)
	servers := tools.MatchedServers()
	for _, server := range servers {
		deployToServer(DeployConfig{
			Config: config.Data, Tasks: server.Tasks,
			Addr:   config.Data.DeployUser + `@` + server.Addr,
			GitTag: tag, GitHost: gitHost,
		})
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
}

var deployTmpl *template.Template

func deployToServer(server DeployConfig) {
	color.Cyan(server.Addr)

	if deployTmpl == nil {
		deployTmpl = template.Must(template.ParseFiles(
			path.Join(config.Root, `config/shell/deploy.tmpl.sh`),
		))
	}

	var buf bytes.Buffer
	err := deployTmpl.Execute(&buf, server)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, server.Addr, buf.String())
}

func getGitHost(gitAddr string) string {
	re := regexp.MustCompile(`@(.*):`)
	return re.FindStringSubmatch(gitAddr)[1]
}
