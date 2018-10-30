package access

import (
	"fmt"

	"github.com/lovego/config/conf"
	"github.com/lovego/xiaomei/release"
)

type Config struct {
	*conf.Conf
	App, Web *service
}

func getConfig(env, svcName, downAddr string) (interface{}, error) {
	if svcName == `` {
		data := Config{
			Conf: release.AppConf(env),
			App:  newService(`app`, env, downAddr),
			Web:  newService(`web`, env, downAddr),
		}
		/*
			if data.App == nil && data.Web == nil {
				return nil, fmt.Error(`neither app nor web service defined.`)
			}
		*/
		return data, nil
	} else {
		data := newService(svcName, env, downAddr)
		if data == nil {
			return nil, fmt.Errorf(`%s service not defined.`, svcName)
		}
		return data, nil
	}
}
