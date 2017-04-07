package deploy

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
)

type Runnable interface {
	FilesForRun() []string
	EnvForRun() []string
	CmdForRun() []string
}

func Run(svcName string, img Runnable) error {
	args := []string{`run`, `-it`, `--rm`,
		`--name=` + release.Name() + `_` + svcName,
	}
	if flags, err := getDriver().FlagsForRun(svcName); err != nil {
		return err
	} else {
		args = append(args, flags...)
	}

	for _, file := range img.FilesForRun() {
		args = append(args, `-v`, file)
	}
	for _, env := range img.EnvForRun() {
		args = append(args, `-e`, env)
	}
	args = append(args, ImageNameOf(svcName))
	if cmd := img.CmdForRun(); cmd != nil {
		args = append(args, cmd...)
	}
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}
