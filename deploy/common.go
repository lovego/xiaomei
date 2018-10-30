package deploy

import (
	"github.com/lovego/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/images"
	//	"github.com/lovego/xiaomei/registry"
)

func GetCommonArgs(svcName, env, timeTag string) []string {
	service := conf.GetService(svcName, env)

	args := []string{}
	if name := images.Get(svcName).EnvironmentEnvVar(); name != `` {
		args = append(args, `-e`, name+`=`+env)
	}
	args = append(args, service.Options...)
	args = append(args, service.ImageNameWithTag(timeTag))
	args = append(args, service.Command...)
	return args
}
