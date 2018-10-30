package new

import (
	"path/filepath"

	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `new <project-path>`,
		Short: `create a new project.`,
		RunE: release.Arg1Call(``, func(dir string) error {
			return New(dir)
		}),
	}
	return cmd
}

func New(proDir string) error {
	config, err := getConfig(proDir)
	if err != nil {
		return err
	}
	tmplsDir := filepath.Join(fs.GetGoSrcPath(), `github.com/lovego/xiaomei/new/webapp`)
	return walk(tmplsDir, proDir, config)
}
