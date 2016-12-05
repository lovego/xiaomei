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

	example := filepath.Join(fmwk.Root(), `example`)
	if !cmd.Ok(cmd.O{}, `cp`, `-r`, example, dir) {
		return
	}

	appName := filepath.Base(proPath)
	script := fmt.Sprintf(`
	cd %s;
	sed -i'' 's/%s/%s/g' main.go
	sed -i'' 's/example/%s/g' .gitignore $(fgrep -rl example release/config)
	`, dir,
		strings.Replace(`/`, `\/`, filepath.Join(fmwk.Path(), `example`), -1),
		strings.Replace(`/`, `\/`, proPath, -1),
		strings.Replace(`/`, `\/`, appName, -1),
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
				fmt.Println(dir, `is not empty.`)
				os.Exit(1)
			}
		} else {
			fmt.Println(dir, `is not a dir.`)
			os.Exit(1)
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
		os.Exit(1)
	}

	if !filepath.IsAbs(dir) {
		var err error
		if dir, err = filepath.Abs(dir); err != nil {
			panic(err)
		}
	}

	goenv := os.Getenv(`GOENV`)
	if goenv == `` {
		fmt.Println(`no GOENV environment variable set.`)
		os.Exit(1)
	}
	goenv = filepath.Join(goenv, `src`)

	rel, err := filepath.Rel(goenv, dir)
	if err != nil {
		panic(err)
	}
	if rel[0] == '.' {
		fmt.Printf("project dir must be under %s.\n", goenv)
		os.Exit(1)
	}
	return rel
}
