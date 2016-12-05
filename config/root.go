package config

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/config/fmwk"
	"github.com/bughou-go/xiaomei/utils"
)

var rootDir string

func Root() string {
	if rootDir != `` {
		return rootDir
	}
	program, cwd := absProgramPath()
	fmwkBin := filepath.Join(os.Getenv(`GOPATH`), `bin`, path.Base(fmwk.Path()))
	if program == fmwkBin /* fmwkBin ... */ ||
		strings.HasSuffix(program, `.test`) /* go test ... */ ||
		strings.HasPrefix(program, `/tmp/`) /* go run ... */ {
		rootDir = filepath.Join(detectRoot(cwd, `release/config/config.yml`), `release`)
	} else {
		// binary under project/release/ dir
		rootDir = detectRoot(filepath.Dir(program), `config/config.yml`)
	}
	return rootDir
}

func absProgramPath() (string, string) {
	program := os.Args[0]
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	switch i := strings.IndexByte(program, filepath.Separator); {
	case i < 0: // path search
		if program, err = exec.LookPath(program); err != nil {
			panic(err)
		}
	case i > 0: // relative path
		program = filepath.Join(cwd, program)
	}
	return program, cwd
}

func detectRoot(dir, feature string) string {
	for ; dir != `/`; dir = filepath.Dir(dir) {
		if utils.Exist(filepath.Join(dir, feature)) {
			return dir
		}
	}
	panic(`project not found.`)
}
