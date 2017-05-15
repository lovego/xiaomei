package tasks

import (
	"fmt"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) Envs() []string {
	return []string{`GOENV=` + release.Env()}
}

func (i Image) BuildDir() string {
	return filepath.Join(release.Root(), `img-app`)
}

func (i Image) Dockerfile() string {
	return `tasksDockerfile`
}

func (i Image) FilesForRun() []string {
	root := filepath.Join(release.Root(), `img-app`)
	name := release.App().Name()
	return []string{
		fmt.Sprintf(`%s/tasks:/home/ubuntu/%s/%s-tasks`, root, name, name),
		fmt.Sprintf(`%s/config:/home/ubuntu/%s/config`, root, name),
		fmt.Sprintf(`%s/views:/home/ubuntu/%s/views`, root, name),
	}
}

func (i Image) EnvsForRun() []string {
	return []string{`GODEV=true`}
}

func (i Image) Prepare() error {
	if err := buildBinary(); err != nil {
		return err
	}
	return nil
}
