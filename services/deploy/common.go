package deploy

import (
	"github.com/lovego/config/config"
	"github.com/lovego/xiaomei/release"
)

func GetCommonArgs(svcName, env, tag string) []string {
	args := []string{`-e`, config.EnvVar + `=` + env}

	service := release.GetService(svcName, env)
	args = append(args, service.Options...)
	args = append(args, service.ImageName(tag))
	args = append(args, service.Command...)
	return args
}
