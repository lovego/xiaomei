package godoc

import (
	"os"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
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
	return filepath.Join(release.Root(), `img-godoc/Dockerfile`)
}
