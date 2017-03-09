package release

import (
	"os"
	"strings"

	"github.com/bughou-go/xiaomei/utils/slice"
)

var theEnv string
var theEnvs []string
var arg1IsEnv bool

func Env() string {
	if theEnv == `` {
		if env := os.Getenv(`GOENV`); env != `` {
			theEnv = env
		} else {
			if strings.HasSuffix(os.Args[0], `.test`) {
				theEnv = `test`
			} else if len(os.Args) >= 2 && slice.ContainsString(Envs(), os.Args[1]) {
				theEnv = os.Args[1]
				arg1IsEnv = true
			} else {
				theEnv = `dev`
			}
		}
	}
	return theEnv
}

func Envs() []string {
	if theEnvs == nil {
		envs := []string{}
		if InProject() {
			for env := range GetClusters() {
				envs = append(envs, env)
			}
		}
		theEnvs = envs
	}
	return theEnvs
}

func Arg1IsEnv() bool {
	Env()
	return arg1IsEnv
}
