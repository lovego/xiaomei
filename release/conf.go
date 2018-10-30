package release

import (
	"path/filepath"

	"github.com/lovego/config/conf"
	"github.com/lovego/fs"
	"github.com/lovego/strmap"
)

var theName string
var appConfig *conf.Config
var appData = map[string]strmap.StrMap{}

func Name() string {
	if theName == `` {
		if AppConfig() != nil {
			theName = AppConfig().Name
		} else {
			theName = filepath.Base(Root())
		}
	}
	return theName
}

func AppConfig() *conf.Config {
	if appConfig == nil {
		if fs.Exist(filepath.Join(Root(), `img-app`)) {
			appConfig = conf.Get(filepath.Join(Root(), `img-app/config/config.yml`))
		} else {
			appConfig = conf.Get(filepath.Join(Root(), `config.yml`))
		}
	}
	return appConfig
}

func AppConf(env string) *conf.Conf {
	return AppConfig().Get(env)
}

func ServiceName(svcName, env string) string {
	return AppConf(env).DeployName() + `-` + svcName
}

func AppData(env string) strmap.StrMap {
	if data := appData[env]; data == nil {
		if fs.Exist(filepath.Join(Root(), `img-app`)) {
			data = conf.Data(filepath.Join(Root(), `img-app/config/envs/`+env+`.yml`))
		} else {
			data = conf.Data(filepath.Join(Root(), `envs/`+env+`.yml`))
		}
		appData[env] = data
		return data
	} else {
		return data
	}
}
