package access

import (
	"fmt"
)

func accessPrint(env, svcName string) error {
	data, err := getConfig(env, svcName)
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
	data, err := getConfig(env, svcName)
	if err != nil {
		return err
	}
	nginxConf, err := getNginxConf(svcName, data)
	if err != nil {
		return err
	}
	return setupNginx(env, feature, nginxConf, data)
}
