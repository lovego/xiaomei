package godoc

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) InstanceEnvName() string {
	return `GODOCPORT`
}

func (i Image) BuildDir() string {
	return filepath.Join(release.Root(), `..`)
}

func (i Image) Dockerfile() string {
	return `release/img-godoc/Dockerfile`
}
