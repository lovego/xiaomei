package deploy

import (
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func ClearDeployTags() error {
	err1 := clearRemoteDeployTags()
	err2 := clearLocalDeployTags()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

// clear remote obsolete deploy tags
func clearRemoteDeployTags() error {
	const historyCount = 5
	const prefix = `refs/tags/`

	lines, err := cmd.Run(cmd.O{Output: true},
		`git`, `ls-remote`, `--tags`, `origin`, prefix+config.Env()+`*`,
	)
	if err != nil {
		return err
	}
	refs := []string{}
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > len(prefix) && isDeployTag(line[len(prefix):]) {
			refs = append(refs, line)
		}
	}
	if len(refs) <= historyCount {
		return nil
	}
	refs = refs[:historyCount]

	args := append([]string{`push`, `--delete`, `origin`}, refs...)

	_, err = cmd.Run(cmd.O{}, `git`, args...)
	return err
}

// clear local obsolete deploy tags
func clearLocalDeployTags() error {
	lines, err := cmd.Run(cmd.O{Output: true},
		`git`, `tag`, `--list`, config.Env()+`*`,
	)
	if err != nil {
		return err
	}
	tags := []string{}
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if isDeployTag(line) {
			tags = append(tags, line)
		}
	}
	for _, tag := range tags {
		if cmd.Fail(cmd.O{NoStdout: true}, `git`, `ls-remote`, `--tags`, `--exit-code`, `origin`, tag) {
			_, err = cmd.Run(cmd.O{}, `git`, `tag`, `-d`, tag)
		}
	}
	return err
}
