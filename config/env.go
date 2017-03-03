package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils/slice"
)

var arg1IsEnv bool

func Arg1IsEnv() bool {
	Env()
	return arg1IsEnv
}

func detectEnv() string {
	env := os.Getenv(`GOENV`)
	if env != `` {
		return env
	}
	if strings.HasSuffix(os.Args[0], `.test`) {
		env = `test`
	} else if os.Args[0] == `xiaomei` && len(os.Args) >= 2 &&
		slice.ContainsString(Envs(), os.Args[1]) {
		env = os.Args[1]
		arg1IsEnv = true
	} else {
		env = `dev`
	}
	return env
}

func availableEnvs() (results []string) {
	if !InProject() {
		return nil
	}
	pathes, err := filepath.Glob(filepath.Join(Root(), `config/envs/*.yml`))
	if err != nil {
		panic(err)
	}
	for _, p := range pathes {
		if env := strings.TrimSuffix(filepath.Base(p), `.yml`); env != `dev` {
			results = append(results, env)
		}
	}
	return
}
