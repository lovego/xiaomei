package release

import (
	"os"
	"path"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils/fs"
)

var theRoot *string

func Root() string {
	root := detectRoot()
	if root == `` {
		panic(`release root not found.`)
	}
	return root
}

func InProject() bool {
	return detectRoot() != ``
}

func detectRoot() string {
	if theRoot == nil {
		if cwd, err := os.Getwd(); err != nil {
			panic(err)
		} else if dir := fs.DetectDir(cwd, `release/stack.yml`); dir != `` {
			dir = filepath.Join(dir, `release`)
			theRoot = &dir
		} else {
			return ``
		}
	}
	return *theRoot
}

// package go path
func Path() string {
	proDir := path.Join(Root(), `../`)

	if !filepath.IsAbs(proDir) {
		var err error
		if proDir, err = filepath.Abs(proDir); err != nil {
			panic(err)
		}
	}

	srcPath, err := fs.GetGoSrcPath()
	if err != nil {
		panic(err)
	}

	proPath, err := filepath.Rel(srcPath, proDir)
	if err != nil {
		panic(err)
	}
	if proPath[0] == '.' {
		panic(`project dir must be under ` + srcPath + "\n")
	}
	return proPath
}