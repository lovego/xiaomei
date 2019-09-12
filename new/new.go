package new

import (
	"path/filepath"

	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var typ string
	cmd := &cobra.Command{
		Use:   `new <project-path>`,
		Short: `create a new project.`,
		RunE: release.Arg1Call(``, func(dir string) error {
			return New(dir, typ)
		}),
	}
	cmd.Flags().StringVarP(&typ, `type`, `t`, `full`, `project type.
 app: only service that provides Golang API.
 web: only service that provides fontend UI.
logc: only service that collect logs to ElasticSearch.
full: all services including app, web and logc.
`)
	return cmd
}

func New(proDir, typ string) error {
	config, err := getConfig(proDir)
	if err != nil {
		return err
	}
	tmplsDir := filepath.Join(fs.GetGoSrcPath(), `github.com/lovego/xiaomei/new/templates/`+typ)
	return walk(tmplsDir, proDir, config)
}
