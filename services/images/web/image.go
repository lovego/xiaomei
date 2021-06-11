package web

import (
	"fmt"
	"path/filepath"

	"github.com/lovego/xiaomei/release"
)

type Image struct {
}

func (i Image) PortEnvVar() string {
	return `NGINXPORT`
}

func (i Image) DefaultPort() uint16 {
	return 80
}

func (i Image) OptionsForRun(env string) []string {
	return []string{
		`-e=SendfileOff=true`,
		fmt.Sprintf("-v=%s/public:/var/www/%s",
			filepath.Join(release.Root(), `img-web`), release.Name(env),
		),
	}
}
