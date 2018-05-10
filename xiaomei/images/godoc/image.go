package godoc

import (
	"os"
	"path/filepath"
)

type Image struct {
}

func (i Image) InstanceEnvName() string {
	return `GODOCPORT`
}

func (i Image) BuildDir() string {
	return filepath.Join(os.Getenv(`GOPATH`), `src`)
}

func (i Image) Dockerfile() string {
	return `release/img-godoc/Dockerfile`
}
