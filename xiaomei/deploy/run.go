package deploy

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
)

func run(svcName string) error {
	image := images.Get(svcName)
	if err := image.PrepareOrBuild(); err != nil {
		return err
	}

	args := []string{`run`, `-it`, `--rm`}
	if flags, err := getDriver().FlagsForRun(svcName); err != nil {
		return err
	} else {
		args = append(args, flags...)
	}
	for _, env := range image.Envs() {
		args = append(args, `-e`, env)
	}
	for _, env := range image.EnvsForRun() {
		args = append(args, `-e`, env)
	}
	for _, file := range image.FilesForRun() {
		args = append(args, `-v`, file)
	}
	for _, file := range conf.VolumesFor(svcName) {
		args = append(args, `-v`, file)
	}
	args = append(args, conf.ImageNameOf(svcName))
	args = append(args, conf.CommandFor(svcName)...)
	_, err := cmd.Run(cmd.O{Print: true}, `docker`, args...)
	return err
}
