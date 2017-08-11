package deploy

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
)

func run(svcName string) error {
	if err := images.Build(svcName, true); err != nil {
		return err
	}

	args := []string{`run`, `-it`, `--rm`}
	if flags, err := getDriver().FlagsForRun(svcName); err != nil {
		return err
	} else {
		args = append(args, flags...)
	}
	image := images.Get(svcName)
	for _, env := range image.Envs() {
		args = append(args, `-e`, env)
	}
	for _, env := range image.EnvsForRun() {
		args = append(args, `-e`, env)
	}
	args = append(args, conf.OptionsFor(svcName)...)
	args = append(args, conf.ImageNameOf(svcName))
	args = append(args, conf.CommandFor(svcName)...)
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}
