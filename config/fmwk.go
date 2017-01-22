package config

import (
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/fs"
)

var Fmwk FmwkConf

type FmwkConf struct {
	root, bin string
}

func (f *FmwkConf) Path() string {
	return `github.com/bughou-go/xiaomei`
}

func (f *FmwkConf) Bin() string {
	if f.bin == `` {
		f.bin, _ = cmd.Run(cmd.O{Output: true}, `which`, `xiaomei`)
	}
	return f.bin
}

func (f *FmwkConf) Root() string {
	if f.root != `` {
		return f.root
	}
	if root := detectRoot(); root != `` {
		if vendorPkg := filepath.Join(root, `vendor`, f.Path()); fs.Exist(vendorPkg) {
			f.root = vendorPkg
			return f.root
		}
	}

	if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, f.Path()); fs.Exist(globalPkg) {
		f.root = globalPkg
	} else {
		panic(`framework not found.`)
	}
	return f.root
}
