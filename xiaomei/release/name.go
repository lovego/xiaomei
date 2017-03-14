package release

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/config/conf"
	"github.com/bughou-go/xiaomei/utils/fs"
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

func App() *conf.Conf {
	if appConf == nil {
		appConf = conf.New(filepath.Join(Root(), `img-app`), Env())
	}
	return appConf
}

func AppIn(env string) *conf.Conf {
	return conf.New(filepath.Join(Root(), `img-app`), env)
}
