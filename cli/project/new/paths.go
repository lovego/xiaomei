package new

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils/fs"
)

const examplePath = `github.com/bughou-go/xiaomei/example`

func getExampleDir() (string, error) {
	srcPath, err := getGoSrcPath()
	if err != nil {
		return ``, err
	}
	exampleDir := filepath.Join(srcPath, examplePath)
	if fs.Exist(exampleDir) {
		return exampleDir, nil
	} else {
		return ``, errors.New(exampleDir + `: not exist.`)
	}
}

func getProjectPath(proDir string) (string, error) {
	if proDir == `` {
		return ``, errors.New(`project path can't be empty.`)
	}

	if !filepath.IsAbs(proDir) {
		var err error
		if proDir, err = filepath.Abs(proDir); err != nil {
			return ``, err
		}
	}

	srcPath, err := getGoSrcPath()
	if err != nil {
		return ``, err
	}

	proPath, err := filepath.Rel(srcPath, proDir)
	if err != nil {
		return ``, err
	}
	if proPath[0] == '.' {
		return ``, errors.New(`project dir must be under ` + srcPath + "\n")
	}
	return proPath, nil
}

func makeProjectDir(dir string) error {
	fi, err := os.Stat(dir)
	switch {
	case err == nil:
		if fi.IsDir() {
			if !fs.IsEmptyDir(dir) {
				return errors.New(dir + ` not empty.`)
			}
		} else {
			return errors.New(dir + ` not dir.`)
		}
	case os.IsNotExist(err):
		return os.MkdirAll(dir, 0775)
	default:
		return err
	}
	return nil
}

func getGoSrcPath() (string, error) {
	gopath := os.Getenv(`GOPATH`)
	if gopath == `` {
		return ``, errors.New(`empty env variable GOPATH.`)
	}
	return filepath.Join(gopath, `src`), nil
}
