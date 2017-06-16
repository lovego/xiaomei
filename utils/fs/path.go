package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func Exist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || !os.IsNotExist(err)
}

func NotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func IsFile(p string) bool {
	fi, _ := os.Stat(p)
	return fi != nil && fi.Mode().IsRegular()
}

func IsDir(p string) bool {
	fi, _ := os.Stat(p)
	return fi != nil && fi.IsDir()
}

func IsEmptyDir(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	} else {
		return false, err
	}
}

func DetectDir(dir string, features ...string) string {
	for ; dir != `/`; dir = filepath.Dir(dir) {
		if hasAllFeatures(dir, features) {
			return dir
		}
	}
	return ``
}

func hasAllFeatures(dir string, features []string) bool {
	for _, feature := range features {
		if !Exist(filepath.Join(dir, feature)) {
			return false
		}
	}
	return true
}

func GetGoSrcPath() (string, error) {
	gopath := os.Getenv(`GOPATH`)
	if gopath == `` {
		return ``, errors.New(`empty env variable GOPATH.`)
	}
	return filepath.Join(gopath, `src`), nil
}
