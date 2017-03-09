package release

import (
	"os"

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
		} else {
			dir := fs.DetectDir(cwd, `release/stack.yml`)
			theRoot = &dir
		}
	}
	return *theRoot
}
