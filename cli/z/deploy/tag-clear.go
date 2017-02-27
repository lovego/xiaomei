package deploy

import (
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/slice"
)

func ClearTags() error {
	if tags, err := ClearRemoteTags(); err == nil && len(tags) > 0 {
		return clearLocalTags(tags)
	} else {
		return err
	}
}

// clear remote obsolete deploy tags
func ClearRemoteTags() ([]string, error) {
	const historyCount = 5

	refs, err := getRemoteTagRefs()
	obsoleteCount := len(refs) - historyCount
	if err != nil || obsoleteCount <= 0 {
		return nil, err
	}
	_, err = cmd.Run(cmd.O{}, `git`, append(
		[]string{`push`, `--delete`, `origin`}, refs[:obsoleteCount]...,
	)...)

	return refs2tags(refs[obsoleteCount:]), err
}

// clear local obsolete deploy tags
func ClearLocalTags() error {
	if refs, err := getRemoteTagRefs(); err != nil {
		return err
	} else {
		return clearLocalTags(refs2tags(refs))
	}
}

const tagRefsPrefix = `refs/tags/`

func getRemoteTagRefs() ([]string, error) {
	lines, err := cmd.Run(cmd.O{Output: true},
		`git`, `ls-remote`, `--tags`, `origin`, tagRefsPrefix+config.App.Env()+`*`,
	)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, line := range strings.Split(lines, "\n") {
		if index := strings.Index(line, tagRefsPrefix); index > 0 {
			line = line[index:]
		} else {
			continue
		}
		line = strings.TrimSpace(line)
		if len(line) > len(tagRefsPrefix) && isDeployTag(line[len(tagRefsPrefix):]) {
			result = append(result, line)
		}
	}
	return result, nil
}

func refs2tags(refs []string) (result []string) {
	for _, ref := range refs {
		result = append(result, ref[len(tagRefsPrefix):])
	}
	return
}

func clearLocalTags(remoteTags []string) error {
	lines, err := cmd.Run(cmd.O{Output: true},
		`git`, `tag`, `--list`, config.App.Env()+`*`,
	)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(lines, "\n") {
		tag := strings.TrimSpace(line)
		if isDeployTag(tag) && !slice.ContainsString(remoteTags, tag) {
			if _, err := cmd.Run(cmd.O{}, `git`, `tag`, `-d`, tag); err != nil {
				return err
			}
		}
	}
	return nil
}
