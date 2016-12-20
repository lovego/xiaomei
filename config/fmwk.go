package config

import (
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
)

var Fmwk fmwk

type fmwk struct {
	root string
}

func (f *fmwk) Path() string {
	return `github.com/bughou-go/xiaomei`
}

func (f *fmwk) Bin() string {
	return filepath.Join(os.Getenv(`GOPATH`), `bin`, filepath.Base(f.Path()))
}

func (f *fmwk) Root() string {
	if f.root != `` {
		return f.root
	}
	if root := detectRoot(); root != `` {
		if vendorPkg := filepath.Join(root, `vendor`, f.Path()); utils.Exist(vendorPkg) {
			f.root = vendorPkg
			return f.root
		}
	}

	if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, f.Path()); utils.Exist(globalPkg) {
		f.root = globalPkg
	} else {
		panic(`framework not found.`)
	}
	return f.root
}
