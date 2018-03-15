package access

import (
	"fmt"

	"github.com/lovego/xiaomei/config/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

type Config struct {
	CheckCert bool
	*conf.Conf
	App, Web *service
}

func getConfig(env, svcName string, checkCert bool) (interface{}, error) {
	if svcName == `` {
		data := Config{
			CheckCert: checkCert,
			Conf:      release.AppConf(env),
			App:       newService(`app`, env, checkCert),
			Web:       newService(`web`, env, checkCert),
		}
		/*
			if data.App == nil && data.Web == nil {
				return nil, fmt.Error(`neither app nor web service defined.`)
			}
		*/
		return data, nil
	} else {
		data := newService(svcName, env, checkCert)
		if data == nil {
			return nil, fmt.Errorf(`%s service not defined.`, svcName)
		}
		return data, nil
	}
}
