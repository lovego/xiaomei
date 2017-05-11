package godoc

import (
	"fmt"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

type Image struct {
}

func (i Image) PortEnvName() string {
	return `GODOCPORT`
}

func (i Image) FilesForRun() []string {
	root := filepath.Join(release.Root(), `..`)
	name := filepath.Base(root)
	return []string{
		fmt.Sprintf(`%s:/home/ubuntu/go/src/%s`, root, name),
	}
}
