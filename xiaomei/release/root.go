package release

import (
	"os"
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
