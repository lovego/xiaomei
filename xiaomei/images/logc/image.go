package logc

import (
	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) Envs() []string {
	return []string{`GOENV=` + release.Env()}
}
