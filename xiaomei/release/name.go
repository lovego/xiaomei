package release

import (
	"path/filepath"

	"github.com/lovego/xiaomei/config/conf"
	"github.com/lovego/xiaomei/utils/fs"
)

var theName string
var appConf *conf.Conf

func Name() string {
	if theName == `` {
		if fs.IsDir(filepath.Join(Root(), `img-app`)) {
			theName = App().Name()
		} else {
			theName = `cluster`
		}
	}
	return theName
}

func DeployName() string {
	return Name() + `_` + Env()
}

func App() *conf.Conf {
	if appConf == nil {
		appConf = conf.New(filepath.Join(Root(), `img-app`), Env())
	}
	return appConf
}

func AppIn(env string) *conf.Conf {
	return conf.New(filepath.Join(Root(), `img-app`), env)
}
