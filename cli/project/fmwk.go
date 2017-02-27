package project

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/fs"
)

const fmwkPath = `github.com/bughou-go/xiaomei`

var fmwkRootDir string

func fmwkRoot() (string, error) {
	if fmwkRootDir != `` {
		return fmwkRootDir, nil
	}
	if root := config.DetectRoot(); root != `` {
		if vendorPkg := filepath.Join(root, `../../vendor`, fmwkPath); fs.Exist(vendorPkg) {
			fmwkRootDir = vendorPkg
			return fmwkRootDir, nil
		}
	}

	if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, fmwkPath); fs.Exist(globalPkg) {
		fmwkRootDir = globalPkg
	} else {
		return ``, errors.New(`framework not found.`)
	}
	return fmwkRootDir, nil
}
