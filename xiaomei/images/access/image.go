package access

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
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
