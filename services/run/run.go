package run

import (
	"fmt"
	"runtime"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/deploy"
	"github.com/lovego/xiaomei/services/deploy/conf"
	"github.com/lovego/xiaomei/services/images"
)

func run(env, svcName string) error {
	image := images.Get(svcName)

	args := []string{
		`run`, `-it`, `--rm`, `--name=` + release.ServiceName(svcName, env) + `.run`,
	}
	if runtime.GOOS == `linux` { // only linux support host network
		args = append(args, `--network=host`)
	}
	if portEnvVar := image.PortEnvVar(); portEnvVar != `` {
		runPort := getRunPort(image, env, svcName)
		args = append(args, `-e`, fmt.Sprintf(`%s=%d`, portEnvVar, runPort))
		if runtime.GOOS != "linux" {
			args = append(args, fmt.Sprintf(`--publish=%d:%d`, runPort, runPort))
		}
	}
	if options := image.OptionsForRun(); len(options) > 0 {
		args = append(args, options...)
	}

	args = append(args, deploy.GetCommonArgs(svcName, env, ``)...)
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}

func getRunPort(image images.Image, env, svcName string) uint16 {
	if ports := conf.GetService(svcName, env).Ports; len(ports) > 0 {
		return ports[0]
	}
	return image.DefaultPort()
}
