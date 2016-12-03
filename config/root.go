package config

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

var projectRootDir, fmwkRootDir string

func Root() string {
	if projectRootDir != `` {
		return projectRootDir
	}
	program, cwd := absProgramPath()
	fmwkBin := path.Base(path.Dir(reflect.TypeOf(struct{}{}).PkgPath()))
	if program == filepath.Join(os.Getenv(`GOPATH`), `bin`, fmwkBin) /* fmwkBin ... */ ||
		strings.HasSuffix(program, `.test`) /* go test ... */ ||
		strings.HasPrefix(program, `/tmp/`) /* go run ... */ {
		projectRootDir = filepath.Join(detectRoot(cwd, `release/config/config.yml`), `release`)
	} else {
		// binary under project/release/ dir
		projectRootDir = detectRoot(filepath.Dir(program), `config/config.yml`)
	}
	return projectRootDir
}

func FmwkRoot() string {
	if fmwkRootDir != `` {
		return fmwkRootDir
	}
	fmwk := path.Dir(reflect.TypeOf(struct{}{}).PkgPath())
	if vendorPkg := filepath.Join(Root(), `vendor`, fmwk); dirExists(vendorPkg) {
		fmwkRootDir = vendorPkg
	} else if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, fmwk); dirExists(globalPkg) {
		fmwkRootDir = globalPkg
	} else {
		panic(`framework not found.`)
	}
	return fmwkRootDir
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
		if dirExists(filepath.Join(dir, feature)) {
			return dir
		}
	}
	panic(`project not found.`)
}

func dirExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
