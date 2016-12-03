package deploy

import (
	"bytes"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"os"
	"path"
	"text/template"
)

type docServerConfig struct {
	Addr                 string
	GOPATH, DeployPath   string
	GitAddr, GitHost     string
	GitBranch, GitCommit string
}

func UpdateDoc(branch, commit string) {
	if err := os.Chdir(config.Root()); err != nil {
		panic(err)
	}

	gitHost := getGitHost(config.Data.GitAddr)
	servers := cli.MatchedServers()
	for _, server := range servers {
		if server.Misc != `doc-server` {
			continue
		}
		gopath := path.Join(config.Data.DeployRoot, `godoc`)
		updateDocServer(docServerConfig{
			Addr:       config.Data.DeployUser + `@` + server.Addr,
			GOPATH:     gopath,
			DeployPath: path.Join(gopath, `src`, config.Data.AppName),
			GitAddr:    config.Data.GitAddr,
			GitHost:    gitHost,
			GitBranch:  branch,
			GitCommit:  commit,
		})
		break
	}
}

func updateDocServer(server docServerConfig) {
	color.Cyan(server.Addr)

	tmpl := template.Must(template.ParseFiles(
		path.Join(config.Root(), `config/shell/doc-server.tmpl.sh`),
	))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, server)
	if err != nil {
		panic(err)
	}
	cmd.Run(cmd.O{Panic: true}, `ssh`, server.Addr, buf.String())
}
