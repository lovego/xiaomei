package develop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/config/fmwk"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func New(dir string) {
	proPath := projectPath(dir)
	checkPkgDir(dir)

	example := filepath.Join(fmwk.Root(), `example`)
	if !cmd.Ok(cmd.O{}, `cp`, `-rT`, example, dir) {
		return
	}

	appName := filepath.Base(proPath)
	script := fmt.Sprintf(`
	cd %s
	sed -i'' 's/example/%s/g' .gitignore $(fgrep -rl example release/config)
	sed -i'' 's/%s/%s/g' main.go
	ln -sf envs/dev.yml release/config/env.yml 2>/dev/null ||
	cp -f release/config/envs/dev.yml release/config/env.yml
	`, dir, appName,
		strings.Replace(filepath.Join(fmwk.Path(), `example`), `/`, `\/`, -1),
		strings.Replace(proPath, `/`, `\/`, -1),
	)
	if !cmd.Ok(cmd.O{}, `sh`, `-c`, script) {
		return
	}
}

func checkPkgDir(dir string) {
	fi, err := os.Stat(dir)
	switch {
	case err == nil:
		if fi.IsDir() {
			if !utils.IsEmptyDir(dir) {
				fmt.Println(dir, `exist and is not empty.`)
				os.Exit(0)
			}
		} else {
			fmt.Println(dir, `exist and is not a dir.`)
			os.Exit(0)
		}
	case os.IsNotExist(err):
		if err := os.MkdirAll(dir, 0775); err != nil {
			panic(err)
		}
	default:
		panic(err)
	}
}

func projectPath(dir string) string {
	if dir == `` {
		fmt.Println(`project dir can't be empty.`)
		os.Exit(0)
	}

	if !filepath.IsAbs(dir) {
		var err error
		if dir, err = filepath.Abs(dir); err != nil {
			panic(err)
		}
	}

	gopath := os.Getenv(`GOPATH`)
	if gopath == `` {
		fmt.Println(`no GOPATH environment variable set.`)
		os.Exit(0)
	}
	gopath = filepath.Join(gopath, `src`)

	rel, err := filepath.Rel(gopath, dir)
	if err != nil {
		panic(err)
	}
	if rel[0] == '.' {
		fmt.Printf("project dir must be under GOPATH(%s).\n", gopath)
		os.Exit(0)
	}
	return rel
}
