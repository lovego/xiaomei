package release

import (
	"os"
	"path/filepath"

	"github.com/lovego/config/config"
	"github.com/lovego/fs"
)

var theName string
var configMap = make(map[string]*config.Config)

func Config(env string) *config.Config {
	if configMap[env] == nil {
		var environ = config.NewEnv(env)
		var envYaml = environ.Minor() + `.yml`

		var file string
		if imgApp := ImageDir(env, `app`); fs.Exist(imgApp) {
			file = filepath.Join(imgApp, `config`, envYaml)
		} else if file = filepath.Join(Root(env), `config`, envYaml); fs.Exist(file) {
		} else {
			file = filepath.Join(Root(env), `config.yml`)
		}
		configMap[env] = config.Get(file, env)
	}
	return configMap[env]
}

func CheckEnv(env string) (string, error) {
	if env == `` {
		env = os.Getenv(`ProENV`)
	}
	if env == `` {
		env = `dev`
	}
	Config(env) // ensure env exist
	return env, nil
}

func Name(env string) string {
	return Config(env).Name
}

func ImageDir(env, svcName string) string {
	return filepath.Join(Root(env), "img-"+svcName)
}

func ServiceName(env, svcName string) string {
	return Config(env).DeployName() + `.` + svcName
}
