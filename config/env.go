package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils/slice"
)

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
	} else {
		env = `dev`
	}
	return env
}

func availableEnvs() (results []string) {
	if !InProject() {
		return nil
	}
	pathes, err := filepath.Glob(filepath.Join(Root(), `*.yml`))
	if err != nil {
		panic(err)
	}
	for _, p := range pathes {
		results = append(results, strings.TrimSuffix(filepath.Base(p), `.yml`))
	}
	return
}
