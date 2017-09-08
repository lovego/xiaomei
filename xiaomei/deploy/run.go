package deploy

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func run(env, svcName string) error {
	if err := images.Build(svcName, env, ``, true); err != nil {
		return err
	}
	image := images.Get(svcName)

	args := []string{
		`run`, `-it`, `--rm`, `--name=` + release.ServiceName(svcName, env) + `.run`,
	}
	if instanceEnvName := image.InstanceEnvName(); instanceEnvName != `` {
		if instances := conf.GetService(svcName, env).Instances(); len(instances) > 0 {
			args = append(args, `-e`, fmt.Sprintf(`%s=%s`, instanceEnvName, instances[0]))
		}
	}
	if runEnvName := image.RunEnvName(); runEnvName != `` {
		args = append(args, `-e`, fmt.Sprintf(`%s=%s`, runEnvName, `true`))
	}

	args = append(args, getCommonArgs(svcName, env, ``)...)
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}
