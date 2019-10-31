package release

import (
	"fmt"
	"os"
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
	data := appData[env]
	if data == nil {
		if fs.Exist(filepath.Join(Root(), `img-app`)) {
			data = conf.Data(filepath.Join(Root(), `img-app/config/envs/`+env+`.yml`))
		} else if fpath := filepath.Join(Root(), `envs/`+env+`.yml`); fs.Exist(fpath) {
			data = conf.Data(fpath)
		}
		appData[env] = data
	}
	return data
}

func CheckEnv(env string) (string, error) {
	if env == `` {
		env = os.Getenv(`GOENV`)
	}
	if env == `` {
		env = `dev`
	}
	if _, ok := AppConfig().Envs[env]; ok {
		return env, nil
	}

	return ``, fmt.Errorf("env %s is not defined in config.yml", env)
}
