package release

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/cmd"
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
			log.Panic(err)
		} else if dir := detectDir(cwd,
			`release/img-app/config/config.yml`,
			`release/img-app/config_*/config.yml`,
			`release/config.yml`,
			`release/config_*.yml`,
		); dir != `` {
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
		Dir:    filepath.Dir(Root()),
	}, GoCmd(), `list`, `.`)
	if err != nil {
		return ``, err
	}
	return strings.TrimSpace(output), nil
}

func detectDir(dir string, features ...string) string {
	for ; dir != `/`; dir = filepath.Dir(dir) {
		if hasAnyFeatures(dir, features) {
			return dir
		}
	}
	return ``
}

func hasAnyFeatures(dir string, features []string) bool {
	for _, feature := range features {
		if out, err := filepath.Glob(filepath.Join(dir, feature)); err != nil {
			log.Panic(err)
		} else if len(out) > 0 {
			return true
		}
	}
	return false
}
