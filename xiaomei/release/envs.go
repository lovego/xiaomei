package release

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/slice"
)

var getEnvs func() []string

func RegisterEnvsGetter(getter func() []string) {
	getEnvs = getter
}

func CheckEnv(env string) (string, error) {
	if env == `` {
		env = `dev`
	}
	if !slice.ContainsString(getEnvs(), env) {
		return ``, fmt.Errorf("env %s not defined in cluster.yml", env)
	}
	return env, nil
}
