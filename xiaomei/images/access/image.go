package access

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PrepareForBuild() error {
	return nil
}

func (i Image) BuildDir() string {
	return filepath.Join(release.Root(), `img-access`)
}

func (i Image) Dockerfile() string {
	return `Dockerfile`
}

func (i Image) EnvsForDeploy() []string {
	return nil
}

func (i Image) FilesForRun() (result []string) {
	if confs, err := filepath.Glob(release.Root() + `img-access/*.conf`); err != nil {
		panic(err)
	} else {
		for _, conf := range confs {
			result = append(result, conf+`:/etc/nginx/sites-enabled/`+filepath.Base(conf))
		}
		return result
	}
}

func (i Image) EnvsForRun() []string {
	return nil
}

func (i Image) CmdForRun() []string {
	return nil
}
