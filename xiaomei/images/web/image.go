package web

import (
	"fmt"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) InstanceEnvName() string {
	return `NGINXPORT`
}

func (i Image) OptionsForRun() []string {
	return []string{
		`-e=SendfileOff=true`,
		fmt.Sprintf("-v=%s/public:/var/www/%s",
			filepath.Join(release.Root(), `img-web`), release.Name(),
		),
	}
}
