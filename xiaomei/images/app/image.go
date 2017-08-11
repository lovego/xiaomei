package app

import (
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
