package config

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

var Deploy deployVar

type deployVar struct {
	conf deployConf
}

type deployConf struct {
	User      string `yaml:"user"`
	Root      string `yaml:"root"`
	GitAddr   string `yaml:"gitAddr"`
	GitBranch string `yaml:"gitBranch"`
}

func (d *deployVar) Name() string {
	return App.Name() + `_` + App.Env()
}
func (d *deployVar) Root() string {
	return d.conf.Root
}
func (d *deployVar) Path() string {
	return filepath.Join(d.Root(), d.Name())
}
func (d *deployVar) User() string {
	return d.conf.User
}

func (d *deployVar) GitAddr() string {
	return d.conf.GitAddr
}

func (d *deployVar) GitBranch() string {
	if d.conf.GitBranch != `` {
		return d.conf.GitBranch
	}
	d.conf.GitBranch, _ = cmd.Run(cmd.O{Output: true, Panic: true},
		`git`, `rev-parse`, `--abbrev-ref`, `HEAD`)
	return d.conf.GitBranch
}
