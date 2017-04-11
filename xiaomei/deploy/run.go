package deploy

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func Run(svcName string) error {
	img := images.Get(svcName)
	if err := img.PrepareOrBuild(conf.ImageNameOf(svcName)); err != nil {
		return err
	}

	args := []string{`run`, `-it`, `--rm`,
		`--name=` + release.Name() + `_` + svcName,
	}
	if flags, err := getDriver().FlagsForRun(svcName); err != nil {
		return err
	} else {
		args = append(args, flags...)
	}

	for _, env := range img.EnvsForRun() {
		args = append(args, `-e`, env)
	}
	for _, file := range img.FilesForRun() {
		args = append(args, `-v`, file)
	}
	args = append(args, conf.ImageNameOf(svcName))
	if cmd := img.CmdForRun(); cmd != nil {
		args = append(args, cmd...)
	}
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}
