package deploy

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func setupDeployTag(commit string) (string, error) {
	branch := config.Deploy.GitBranch()
	// 拉取最新代码及tag
	o := cmd.O{}
	if _, err := cmd.Run(o, `git`, `fetch`, `origin`, branch, `--tags`, `--prune`); err != nil {
		return ``, err
	}
	tag, err := getDeployTag(commit)
	if err != nil {
		return ``, err
	}
	// 推送最新代码及tag
	if _, err := cmd.Run(o, `git`, `push`, `origin`, branch, tag); err != nil {
		return ``, err
	}
	return tag, nil
}

func getDeployTag(commit string) (string, error) {
	branch := config.Deploy.GitBranch()

	commit = strings.TrimSpace(commit)
	if commit == `` {
		commit = branch
	} else {
		if err := checkDeployCommit(commit, branch); err != nil {
			return ``, err
		}
		if isDeployTag(commit) {
			return commit, nil
		}
	}

	if tag, err := existDeployTag(commit); tag != `` || err != nil {
		return tag, err
	}
	return newDeployTag(commit)
}

// 检查commit
func checkDeployCommit(commit, branch string) error {
	var o = cmd.O{NoStdout: true, NoStderr: true}
	// commit必须存在
	if !cmd.Ok(o, `git`, `rev-parse`, commit) {
		return errors.New(`commit ` + commit + ` not exists.`)
	}
	// commit必须在指定的分支上
	if !cmd.Ok(o, `git`, `merge-base`, `--is-ancestor`, commit, branch) &&
		!cmd.Ok(o, `git`, `merge-base`, `--is-ancestor`, commit, `remotes/origin/`+branch) {
		return errors.New(`commit ` + commit + ` not on branch ` + branch + `.`)
	}
	return nil
}

const tagTimeLayout = `0102-150405`

var deployTagTime = regexp.MustCompile(`^\d{4}-\d{6}$`)

func isDeployTag(tag string) bool {
	env := config.App.Env()
	l := len(env)
	return len(tag) == l+len(tagTimeLayout) && tag[:l] == env && deployTagTime.MatchString(tag[l:])
}

// 已有的deployTag
func existDeployTag(commit string) (string, error) {
	output, err := cmd.Run(cmd.O{Output: true}, `git`, `tag`, `--points-at`, commit)
	if err != nil {
		return ``, err
	}
	if output != `` {
		for _, tag := range strings.Split(output, "\n") {
			tag = strings.TrimSpace(tag)
			if isDeployTag(tag) {
				return tag, nil
			}
		}
	}
	return ``, nil
}

// 新建deployTag
func newDeployTag(commit string) (string, error) {
	tag := config.App.Env() + time.Now().Format(tagTimeLayout)
	if _, err := cmd.Run(cmd.O{}, `git`, `tag`, tag, commit); err != nil {
		return ``, err
	}
	return tag, nil
}
