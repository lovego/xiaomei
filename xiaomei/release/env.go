package release

import (
	"os"
	"strings"

	"github.com/lovego/xiaomei/utils/slice"
)

var theEnv string
var arg1IsEnv bool

func Env() string {
	if theEnv == `` {
		if env := os.Getenv(`GOENV`); env != `` {
			theEnv = env
		} else {
			if strings.HasSuffix(os.Args[0], `.test`) {
				theEnv = `test`
			} else if len(os.Args) >= 2 && slice.ContainsString(getEnvs(), os.Args[1]) {
				theEnv = os.Args[1]
				arg1IsEnv = true
			} else {
				theEnv = `dev`
			}
		}
	}
	return theEnv
}

func Arg1IsEnv() bool {
	Env()
	return arg1IsEnv
}

var getEnvs func() []string

func RegisterEnvsGetter(getter func() []string) {
	getEnvs = getter
}
