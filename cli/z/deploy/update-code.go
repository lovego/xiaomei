package deploy

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func updateCode(sshAddr, tag string) {
	if config.IsLocalEnv() {
		cmd.Run(cmd.O{Panic: true}, `sh`, `-c`, getUpdatingCodeScript(tag))
	} else {
		cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, sshAddr, getUpdatingCodeScript(tag))
	}
}

var updatingCodeScript string

func getUpdatingCodeScript(tag string) string {
	if updatingCodeScript == `` {
		var buf bytes.Buffer
		if err := template.Must(template.New(``).Parse(updateCodeTmpl)).
			Execute(&buf, updatingCodeConfig(tag)); err != nil {
			panic(err)
		}
		updatingCodeScript = buf.String()
	}

	return updatingCodeScript
}

func updatingCodeConfig(tag string) interface{} {
	gitAddr := config.Deploy.GitAddr()
	if gitAddr == `` {
		panic(`empty config gitAddr.`)
	}
	return struct {
		config.Conf
		GitTag, GitHost, GitAddr string
	}{
		Conf:    config.Data(),
		GitTag:  tag,
		GitHost: getGitHost(gitAddr),
		GitAddr: gitAddr,
	}
}

func getGitHost(gitAddr string) string {
	re := regexp.MustCompile(`@(.*):`)
	return re.FindStringSubmatch(gitAddr)[1]
}

const updateCodeTmpl = `
main () {
  clone_code_if_should
  cd {{.Deploy.Path}}/release || exit 1
  reset_code_to_target_point
  mkdir -p bins
}

clone_code_if_should() {
  test -d {{.Deploy.Path}} ||
    sudo mkdir -p {{.Deploy.Path}} &&
    sudo chown -R $(id -un):$(id -gn) {{.Deploy.Path}} || exit 1

  if test ! -d {{.Deploy.Path}}/.git; then
    ssh-keygen -F {{.GitHost}} > /dev/null || ssh-keyscan -H {{.GitHost}} >> ~/.ssh/known_hosts
    git clone --depth=1 {{.Deploy.GitAddr}} {{.Deploy.Path}} || exit 1
  fi
}

reset_code_to_target_point() {
  git fetch origin -u --tags {{.Deploy.GitBranch}}:{{.Deploy.GitBranch}} || exit 1
  git checkout -q {{.Deploy.GitBranch}} || exit 1

  if git merge-base --is-ancestor {{.GitTag}} {{.Deploy.GitBranch}}; then
    git reset --hard {{.GitTag}} || exit 1
  else
    echo tag {{.GitTag}} is not in branch {{.Deploy.GitBranch}}; exit 1
  fi
}

main
`
