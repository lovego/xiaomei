package deploy

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path"
	"regexp"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"text/template"
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

	git_host := getGitHost(config.Data.GitAddr)
	servers := tools.GetMatchedServers()
	for _, server := range servers {
		deployToServer(DeployConfig{
			Config: config.Data, Tasks: server.Tasks,
			Addr:   config.Data.DeployUser + `@` + server.Addr,
			GitTag: tag, GitHost: git_host,
		})
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
}

var deploy_tmpl *template.Template

func deployToServer(server DeployConfig) {
	color.Cyan(server.Addr)

	if deploy_tmpl == nil {
		deploy_tmpl = template.Must(template.ParseFiles(
			path.Join(config.Root, `config/shell/deploy.tmpl.sh`),
		))
	}

	var buf bytes.Buffer
	err := deploy_tmpl.Execute(&buf, server)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, server.Addr, buf.String())
}

func getGitHost(git_addr string) string {
	re := regexp.MustCompile(`@(.*):`)
	return re.FindStringSubmatch(git_addr)[1]
}
