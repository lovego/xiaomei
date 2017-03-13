package access

import (
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PrepareForBuild() error {
	return nil
}

func (i Image) BuildDir() string {
	return release.Root()
}

func (i Image) Dockerfile() string {
	return `Dockerfile`
}

func (i Image) FilesForRun() []string {
	return []string{
		release.App().Root() + `/sites:/etc/nginx/sites-enabled`,
	}
}

func (i Image) EnvForRun() []string {
	return nil
}

func (i Image) CmdForRun() []string {
	return []string{`sh`, `-c`, `nginx -t && nginx`}
}
