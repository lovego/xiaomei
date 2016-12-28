package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils"
)

func InProject() bool {
	return detectRoot() != ``
}

func detectRoot() string {
	program, cwd := absProgramPath()
	if program == Fmwk.Bin() /* fmwk ... */ ||
		strings.HasSuffix(program, `.test`) /* go test ... */ ||
		strings.HasPrefix(program, `/tmp/`) /* go run ... */ {
		if dir := detectDir(cwd, `release/config/config.yml`); dir == `` {
			return ``
		} else {
			return filepath.Join(dir, `release`)
		}
	} else {
		// binary under project/release/ dir
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
		if utils.Exist(filepath.Join(dir, feature)) {
			return dir
		}
	}
	return ``
}
