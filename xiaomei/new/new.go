package new

import (
	"path/filepath"

	"github.com/lovego/xiaomei/utils/fs"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var isInfra bool
	cmd := &cobra.Command{
		Use:   `new <project-path>`,
		Short: `create a new project.`,
		RunE: release.Arg1Call(``, func(dir string) error {
			return New(dir, isInfra)
		}),
	}
	cmd.Flags().BoolVarP(&isInfra, `infra`, `i`, false, `create a infrastructures project.`)
	return cmd
}

func New(proDir string, isInfra bool) error {
	config, err := getConfig(proDir)
	if err != nil {
		return err
	}
	tmplsDir, err := getTemplatesDir(isInfra)
	if err != nil {
		return err
	}
	return walk(tmplsDir, proDir, config)
}

func getTemplatesDir(isInfra bool) (string, error) {
	srcPath, err := fs.GetGoSrcPath()
	if err != nil {
		return ``, err
	}
	tmplsDir := filepath.Join(srcPath, `github.com/lovego/xiaomei/xiaomei/new`)
	if isInfra {
		return filepath.Join(tmplsDir, `infra`), nil
	} else {
		return filepath.Join(tmplsDir, `webapp`), nil
	}
}
