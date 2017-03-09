package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/config/conf"
	"github.com/bughou-go/xiaomei/utils/fs"
)

var Config = getConfig()

func getConfig() conf.Conf {
	root := detectRoot()
	if root == `` {
		panic(`app root not found.`)
	}
	return conf.New(root, detectEnv())
}

func detectRoot() string {
	program, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(program, `.test`) /* go test ... */ ||
		strings.HasPrefix(program, `/tmp/`) /* go run ... */ {
		if cwd, err := os.Getwd(); err != nil {
			panic(err)
		} else if dir := fs.DetectDir(cwd, `release/stack.yml`); dir == `` {
			return ``
		} else {
			return filepath.Join(dir, `release/img-app`)
		}
	} else { // project binary file
		return fs.DetectDir(filepath.Dir(program), `config/config.yml`)
	}
}

func detectEnv() string {
	env := os.Getenv(`GOENV`)
	if env != `` {
		return env
	}
	if strings.HasSuffix(os.Args[0], `.test`) {
		env = `test`
	} else {
		env = `dev`
	}
	return env
}
