package fmwk

import (
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
)

var rootDir string

func Path() string {
	return `github.com/bughou-go/xiaomei`
}

func Root() string {
	if rootDir != `` {
		return rootDir
	}
	/*
		if vendorPkg := filepath.Join(Root(), `vendor`, Path()); utils.Exist(vendorPkg) {
			rootDir = vendorPkg
		} else
	*/
	if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, Path()); utils.Exist(globalPkg) {
		rootDir = globalPkg
	} else {
		panic(`framework not found.`)
	}
	return rootDir
}
