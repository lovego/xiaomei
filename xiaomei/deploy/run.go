package deploy

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func run(svcName string) error {
	if err := images.Build(svcName, true); err != nil {
		return err
	}

	args := []string{
		`run`, `-it`, `--rm`, `--network=host`,
		`--name=` + release.DeployName() + `_` + svcName + `.run`,
	}
	image := images.Get(svcName)
	if instanceEnvName := image.InstanceEnvName(); instanceEnvName != `` {
		if instances := conf.InstancesOf(svcName); len(instances) > 0 {
			args = append(args, fmt.Sprintf(`-e=%s=%s`, instanceEnvName, instances[0]))
		}
	}
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
