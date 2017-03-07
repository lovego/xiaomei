package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/fs"
)

func InProject() bool {
	return DetectRoot() != ``
}

var theRootDir interface{}

func DetectRoot() string {
	if theRootDir == nil {
		theRootDir = detectRoot()
	}
	return theRootDir.(string)
}

func detectRoot() string {
	program, cwd := absProgramPath()
	fmwkBin, _ := cmd.Run(cmd.O{Output: true}, `which`, `xiaomei`)
	if program == fmwkBin /* fmwk ... */ ||
		strings.HasSuffix(program, `.test`) /* go test ... */ ||
		strings.HasPrefix(program, `/tmp/`) /* go run ... */ {
		if dir := detectDir(cwd, `release/stack.yml`); dir == `` {
			return ``
		} else {
			return filepath.Join(dir, `release/img-app`)
		}
	} else {
		// project binary file
		return detectDir(filepath.Dir(program), `config/config.yml`)
	}
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

func detectDir(dir, feature string) string {
	for ; dir != `/`; dir = filepath.Dir(dir) {
		if fs.Exist(filepath.Join(dir, feature)) {
			return dir
		}
	}
	return ``
}
