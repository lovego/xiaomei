package app

import (
	"fmt"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PortEnvName() string {
	return `GOPORT`
}

func (i Image) Envs() []string {
	return []string{`GOENV=` + release.Env()}
}

func (i Image) FilesForRun() []string {
	root := release.App().Root()
	name := release.App().Name()
	return []string{
		fmt.Sprintf(`%s/%s:/home/ubuntu/%s/%s`, root, name, name, name),
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
	/*
		if err :=	Assets(nil); err != nil {
			return err
		}
	*/
	return nil
}
