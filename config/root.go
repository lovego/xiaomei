package config

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bughou-go/xiaomei/utils"
)

var projectRootDir, fmwkRootDir string

func Root() string {
	if projectRootDir != `` {
		return projectRootDir
	}
	program, cwd := absProgramPath()
	fmwkBin := path.Base(FmwkPath())
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

func FmwkPath() string {
	return path.Dir(reflect.TypeOf(struct{}{}).PkgPath())
}

func FmwkRoot() string {
	if fmwkRootDir != `` {
		return fmwkRootDir
	}
	fmwk := FmwkPath()
	if vendorPkg := filepath.Join(Root(), `vendor`, fmwk); utils.Exist(vendorPkg) {
		fmwkRootDir = vendorPkg
	} else if globalPkg := filepath.Join(os.Getenv(`GOPATH`), `src`, fmwk); utils.Exist(globalPkg) {
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
		if utils.Exist(filepath.Join(dir, feature)) {
			return dir
		}
	}
	panic(`project not found.`)
}
