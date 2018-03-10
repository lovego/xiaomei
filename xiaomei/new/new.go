package new

import (
	"path/filepath"

	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/xiaomei/release"
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
	tmplsDir, err := getTemplatesDir()
	if err != nil {
		return err
	}
	return walk(tmplsDir, proDir, config)
}

func getTemplatesDir() (string, error) {
	srcPath, err := fs.GetGoSrcPath()
	if err != nil {
		return ``, err
	}
	tmplsDir := filepath.Join(srcPath, `github.com/lovego/xiaomei/xiaomei/new`)
	return filepath.Join(tmplsDir, `webapp`), nil
}
