package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	//	"github.com/lovego/xiaomei/xiaomei/registry"
)

func getCommonArgs(env, svcName, timeTag string) []string {
	service := conf.GetService(env, svcName)

	args := []string{`--network=host`}
	if name := images.Get(svcName).EnvironmentEnvName(); name != `` {
		args = append(args, `-e`, name+`=`+env)
	}
	args = append(args, service.Options...)
	args = append(args, service.ImageNameWithTag(timeTag))
	args = append(args, service.Command...)
	return args
}
