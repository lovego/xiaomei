package fmwk

import (
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/bughou-go/xiaomei/utils"
)

var rootDir string

func Path() string {
	return path.Dir(path.Dir(reflect.TypeOf(struct{}{}).PkgPath()))
}

func Root() string {
	if rootDir != `` {
		return rootDir
	}
	fmwk := Path()
	if vendorPkg := filepath.Join(Root(), `vendor`, fmwk); utils.Exist(vendorPkg) {
		rootDir = vendorPkg
	} else if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, fmwk); utils.Exist(globalPkg) {
		rootDir = globalPkg
	} else {
		panic(`framework not found.`)
	}
	return rootDir
}
