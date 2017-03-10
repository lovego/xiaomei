package web

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PrepareForBuild() error {
	return nil
}

func (i Image) BuildDir() string {
	return filepath.Join(release.Root(), `img-web`)
}

func (i Image) Dockerfile() string {
	return `Dockerfile`
}

func (i Image) RunFiles() []string {
	return []string{
		filepath.Join(release.Root(), `img-web/public`) + `:/var/www/` + release.Name(),
	}
}
