package release

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lovego/config/config"
	"github.com/lovego/fs"
	"github.com/lovego/strmap"
)

const EnvironmentEnvVar = "ProENV"

var theName string
var configMap = make(map[string]*config.Config)
var dataMap = make(map[string]strmap.StrMap)

func Config(envStr string) *config.Config {
	env := config.NewEnv(envStr)
	if configMap[env.Major()] == nil {
		if fs.Exist(filepath.Join(Root(), `img-app`)) {
			configMap[env.Major()] = config.Get(filepath.Join(
				Root(), `img-app`, env.ConfigDir(), `config.yml`,
			), env.Major())
		} else {
			configMap[env.Major()] = config.Get(filepath.Join(
				Root(), env.ConfigDir()+`.yml`,
			), env.Major())
		}
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

func CheckEnv(env string) (string, error) {
	if env == `` {
		env = os.Getenv(`ProENV`)
	}
	if env == `` {
		env = `dev`
	}
	if _, ok := Config(env).Envs[env]; ok {
		return env, nil
	}

	return ``, fmt.Errorf("env %s is not defined in config.yml", env)
}

func EnvData(envStr string) strmap.StrMap {
	data := dataMap[envStr]
	if data == nil {
		env := config.NewEnv(envStr)
		if fs.Exist(filepath.Join(Root(), `img-app`)) {
			data = config.Data(filepath.Join(
				Root(), `img-app`, env.ConfigDir(), `envs`, env.Minor()+`.yml`,
			))
		} else {
			data = strmap.StrMap{}
		}
		dataMap[envStr] = data
	}
	return data
}
