package project

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/fs"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `new <project-name>`,
		Short: `create a new project.`,
		RunE: func(c *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return errors.New(`<project-name> is required.`)
			case 1:
				return New(args[0])
			default:
				return errors.New(`redundant args.`)
			}
		},
	}
}

func New(proDir string) error {
	proPath, err := projectPath(proDir)
	if err != nil {
		return err
	}
	if err = checkPkgDir(proDir); err != nil {
		return err
	}

	var exampleDir string
	if fmwkRootDir, err := fmwkRoot(); err != nil {
		return err
	} else {
		exampleDir = filepath.Join(fmwkRootDir, `example`)
		if !cmd.Ok(cmd.O{}, `cp`, `-rT`, exampleDir, proDir) {
			return errors.New(`cp templates failed.`)
		}
	}

	proName := filepath.Base(proPath)
	script := fmt.Sprintf(`
	cd %s
	sed -i'' 's/example/%s/g' .gitignore $(fgrep -rl example release/config)
	sed -i'' 's/%s/%s/g' main.go
	sed -i'' 's/secret-string/%s/g' release/config/envs/production.yml
	`, proDir, proName,
		strings.Replace(filepath.Join(fmwkPath, `example`), `/`, `\/`, -1),
		strings.Replace(proPath, `/`, `\/`, -1),
		generateSecret(),
	)
	if !cmd.Ok(cmd.O{}, `sh`, `-c`, script) {
		return errors.New(`process templates failed.`)
	}
	return nil
}

// 32 byte hex string
func generateSecret() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func checkPkgDir(dir string) error {
	fi, err := os.Stat(dir)
	switch {
	case err == nil:
		if fi.IsDir() {
			if !fs.IsEmptyDir(dir) {
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
