package release

import (
	"os"
	"path/filepath"

	"github.com/lovego/config/config"
	"github.com/lovego/fs"
	"github.com/lovego/strmap"
)

var theName string
var configMap = make(map[string]*config.Config)
var dataMap = make(map[string]strmap.StrMap)

func Config(envStr string) *config.Config {
	env := config.NewEnv(envStr)
	if configMap[env.Major()] == nil {
		var file string
		if imgApp := ServiceDir(`app`); fs.Exist(imgApp) {
			file = filepath.Join(imgApp, env.ConfigDir(), `config.yml`)
		} else {
			file = filepath.Join(Root(), env.ConfigDir()+`.yml`)
		}
		configMap[env.Major()] = config.Get(file, env.Major())
	}
	return configMap[env.Major()]
}

func Name(env string) string {
	return Config(env).Name
}

func EnvConfig(env string) *config.EnvConfig {
	return Config(env).Get(env)
}

func ServiceName(svcName, env string) string {
	return EnvConfig(env).DeployName() + `.` + svcName
}

func ServiceDir(svcName string) string {
	return filepath.Join(Root(), "img-"+svcName)
}

func CheckEnv(env string) (string, error) {
	if env == `` {
		env = os.Getenv(`ProENV`)
	}
	if env == `` {
		env = `dev`
	}
	Config(env).Get(env) // ensure env exist
	return env, nil
}

func EnvData(envStr string) strmap.StrMap {
	data := dataMap[envStr]
	if data == nil {
		env := config.NewEnv(envStr)
		var dir string
		if imgApp := ServiceDir(`app`); fs.Exist(imgApp) {
			dir = filepath.Join(imgApp, env.ConfigDir(), `envs`)
		} else {
			dir = filepath.Join(Root(), env.TailMajor(`envs`))
		}
		file := filepath.Join(dir, env.Minor()+`.yml`)
		if fs.Exist(file) {
			data = config.Data(file)
			dataMap[envStr] = data
		}
	}
	return data
}
