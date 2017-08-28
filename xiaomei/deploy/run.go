package deploy

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func run(env, svcName string) error {
	if err := images.Build(env, svcName, true); err != nil {
		return err
	}
	image := images.Get(svcName)
	service := conf.GetService(env, svcName)

	args := []string{
		`run`, `-it`, `--rm`, `--name=` + release.ServiceName(env, svcName) + `.run`,
	}
	if instanceEnvName := image.InstanceEnvName(); instanceEnvName != `` {
		if instances := conf.InstancesOf(env, svcName); len(instances) > 0 {
			args = append(args, `-e`, fmt.Sprintf(`%s=%s`, instanceEnvName, instances[0]))
		}
	}
	if runEnvName := image.RunEnvName(); runEnvName != `` {
		args = append(args, `-e`, fmt.Sprintf(`%s=%s`, runEnvName, `true`))
	}

	args = append(args, getCommonArgs(service, image, env)...)
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}

func getCommonArgs(service conf.Service, image images.Image, env string) []string {
	args := []string{`--network=host`}
	if name := image.EnvironmentEnvName(); name != `` {
		args = append(args, `-e`, name+`=`+env)
	}
	args = append(args, service.Options...)
	args = append(args, image.NameWithDigestInRegistry(env))
	args = append(args, service.Command...)
	return args
}
