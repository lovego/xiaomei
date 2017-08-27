package release

import (
	"path/filepath"

	"github.com/lovego/xiaomei/config/conf"
	"github.com/lovego/xiaomei/utils/fs"
	"github.com/lovego/xiaomei/utils/strmap"
)

var theName string
var appConfig *conf.Config
var appData = map[string]strmap.StrMap{}

func Name() string {
	if theName == `` {
		if fs.IsDir(filepath.Join(Root(), `img-app`)) {
			theName = AppConfig().Name
		} else {
			theName = filepath.Base(Root())
		}
	}
	return theName
}

func AppConfig() *conf.Config {
	if appConfig == nil {
		appConfig = conf.Get(filepath.Join(Root(), `img-app`))
	}
	return appConfig
}

func AppConf(env string) *conf.Conf {
	return AppConfig().Get(env)
}

func ServiceName(env, svcName string) string {
	return AppConf(env).DeployName() + `_` + svcName
}

func AppData(env string) strmap.StrMap {
	if data := appData[env]; data == nil {
		data = conf.Data(filepath.Join(Root(), `img-app`), env)
		appData[env] = data
		return data
	} else {
		return data
	}
}
