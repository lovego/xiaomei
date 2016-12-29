package develop

import (
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func New(dir string) error {
	proPath, err := projectPath(dir)
	if err != nil {
		return err
	}
	if err = checkPkgDir(dir); err != nil {
		return err
	}

	example := filepath.Join(config.Fmwk.Root(), `example`)
	if !cmd.Ok(cmd.O{}, `cp`, `-rT`, example, dir) {
		return errors.New(`cp templates failed.`)
	}

	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	port := r.Intn(10000) + 30000
	appName := filepath.Base(proPath)
	script := fmt.Sprintf(`
	cd %s
	sed -i'' 's/example/%s/g' .gitignore $(fgrep -rl example release/config)
	sed -i'' 's/%s/%s/g' main.go
	sed -i'' 's/secret-string/%s/g' release/config/envs/production.yml
	sed -i'' 's/3000/%d/g' release/config/config.yml
	ln -sf envs/dev.yml release/config/env.yml 2>/dev/null ||
	cp -f release/config/envs/dev.yml release/config/env.yml
	`, dir, appName,
		strings.Replace(filepath.Join(config.Fmwk.Path(), `example`), `/`, `\/`, -1),
		strings.Replace(proPath, `/`, `\/`, -1),
		generateSecret(), port,
	)
	if !cmd.Ok(cmd.O{}, `sh`, `-c`, script) {
		return errors.New(`process templates failed.`)
	}
	return nil
}

// 32 byte hex string
func generateSecret() string {
	b := make([]byte, 16)
	if _, err := crand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func checkPkgDir(dir string) error {
	fi, err := os.Stat(dir)
	switch {
	case err == nil:
		if fi.IsDir() {
			if !utils.IsEmptyDir(dir) {
				return errors.New(dir + ` exist and is not empty.`)
			}
		} else {
			return errors.New(dir + ` exist and is not a dir.`)
		}
	case os.IsNotExist(err):
		if err := os.MkdirAll(dir, 0775); err != nil {
			panic(err)
		}
	default:
		panic(err)
	}
	return nil
}

func projectPath(dir string) (string, error) {
	if dir == `` {
		return ``, errors.New(`project name can't be empty.`)
	}

	if !filepath.IsAbs(dir) {
		var err error
		if dir, err = filepath.Abs(dir); err != nil {
			panic(err)
		}
	}

	gopath := os.Getenv(`GOPATH`)
	if gopath == `` {
		return ``, errors.New(`no GOPATH environment variable set.`)

	}
	gopath = filepath.Join(gopath, `src`)

	rel, err := filepath.Rel(gopath, dir)
	if err != nil {
		panic(err)
	}
	if rel[0] == '.' {
		return ``, errors.New(`project dir must be under GOPATH(` + gopath + ").\n")
	}
	return rel, nil
}
