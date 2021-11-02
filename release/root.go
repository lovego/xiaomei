package release

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/config/config"
	"github.com/lovego/fs"
)

// Root returns an environment's root directory of release config.
func Root(env string) string {
	root := getRoot(env)
	if root == `` {
		log.Fatalf(`release root of %s not found.`, env)
	}
	return root
}

func InProject(env string) bool {
	return getRoot(env) != ""
}

func configFile(env, file string) string {
	return filepath.Join(Root(env), file)
}

var roots = make(map[string]string)

func getRoot(env string) string {
	environ := config.NewEnv(env)
	root, ok := roots[environ.Major()]
	if !ok {
		root = config.DetectReleaseConfigDirOf(environ.Major())
		roots[environ.Major()] = root
	}
	return root
}

var projectRoot *string

func ProjectRoot() string {
	if projectRoot == nil {
		if releaseDir, _ := config.DetectReleaseDir(); releaseDir == `` {
			log.Fatal(`release root not found.`)
		} else {
			projectDir := filepath.Dir(releaseDir)
			projectRoot = &projectDir
		}
	}
	return *projectRoot
}

func ModulePath() (string, error) {
	output, err := cmd.Run(cmd.O{
		Output: true,
		Dir:    ProjectRoot(),
	}, GoCmd(), `list`, `.`)
	if err != nil {
		return ``, err
	}
	return strings.TrimSpace(output), nil
}

func SrcDir() string {
	dir := ProjectRoot()
	if fs.Exist(filepath.Join(dir, "src", "go.mod")) {
		return filepath.Join(dir, "src")
	}
	return dir
}
