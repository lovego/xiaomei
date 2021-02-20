package release

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/fs"
)

var theRoot *string

func Root() string {
	root := detectRoot()
	if root == `` {
		log.Fatal(`release root not found.`)
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
		} else if dir := fs.DetectDir(cwd, `release/img-app/config/config.yml`); dir != `` {
			dir = filepath.Join(dir, `release`)
			theRoot = &dir
		} else if dir := fs.DetectDir(cwd, `release/config.yml`); dir != `` {
			dir = filepath.Join(dir, `release`)
			theRoot = &dir
		} else {
			return ``
		}
	}
	return *theRoot
}

func ModulePath() (string, error) {
	output, err := cmd.Run(cmd.O{
		Output: true,
		Dir:    Root(),
	}, `go`, `list`, `.`)
	if err != nil {
		return ``, err
	}
	return strings.TrimSpace(output), nil
}
