package pro

import (
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/fs"
)

const fmwkPath = `github.com/bughou-go/xiaomei`

var fmwkRootDir string

func fmwkRoot() string {
	if fmwkRootDir != `` {
		return fmwkRootDir
	}
	if root := config.DetectRoot(); root != `` {
		if vendorPkg := filepath.Join(root, `../../vendor`, fmwkPath); fs.Exist(vendorPkg) {
			fmwkRootDir = vendorPkg
			return fmwkRootDir
		}
	}

	if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, fmwkPath); fs.Exist(globalPkg) {
		fmwkRootDir = globalPkg
	} else {
		panic(`framework not found.`)
	}
	return fmwkRootDir
}
