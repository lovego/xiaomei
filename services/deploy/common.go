package deploy

import (
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
	//	"github.com/lovego/xiaomei/registry"
)

func GetCommonArgs(svcName, env, tag string) []string {
	service := release.GetService(svcName, env)

	args := []string{}
	if name := images.Get(svcName).EnvironmentEnvVar(); name != `` {
		args = append(args, `-e`, name+`=`+env)
	}
	args = append(args, service.Options...)
	args = append(args, service.ImageName(tag))
	args = append(args, service.Command...)
	return args
}
