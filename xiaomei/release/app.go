package release

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/config/conf"
	"github.com/bughou-go/xiaomei/utils/fs"
)

var appConf conf.Conf
var theName string

func App() *conf.Conf {
	return conf.New(filepath.Join(Root(), `image-app`), Env())
}

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
