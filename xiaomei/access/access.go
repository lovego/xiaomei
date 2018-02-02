package access

import (
	"fmt"

	"github.com/lovego/xiaomei/xiaomei/release"
)

func accessPrint(env, svcName string) error {
	data, err := getConfig(env, svcName, false)
	if err != nil {
		return err
	}
	nginxConf, err := getNginxConf(svcName, data)
	if err != nil {
		return err
	}
	fmt.Print(nginxConf)
	return nil
}

func accessSetup(env, svcName, feature string) error {
	if release.AppConf(env).Https {
		if err := setupNginx(env, svcName, feature, true); err != nil {
			return err
		}
	}
	return setupNginx(env, svcName, feature, false)
}
