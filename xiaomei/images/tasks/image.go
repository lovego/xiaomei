package tasks

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) EnvironmentEnvName() string {
	return `GOENV`
}

func (i Image) RunEnvName() string {
	return `GODEV`
}

func (i Image) BuildDir() string {
	return filepath.Join(release.Root(), `img-app`)
}

func (i Image) Dockerfile() string {
	return `tasksDockerfile`
}

func (i Image) Prepare() error {
	if err := buildBinary(); err != nil {
		return err
	}
	return nil
}
