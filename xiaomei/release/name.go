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
		if fs.IsDir(filepath.Join(Root(), `image-app`)) {
			theName = App().Name()
		} else {
			theName = `cluster`
		}
	}
	return theName
}

func App() *conf.Conf {
	if appConf == nil {
		appConf = conf.New(filepath.Join(Root(), `image-app`), Env())
	}
	return appConf
}
