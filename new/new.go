package new

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/misc/utils"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var typ, domain string
	var force bool
	cmd := &cobra.Command{
		Use: `new <project path> <registry> [flags]

project path: Go module path or directory path for the project, required. The last element of project path is used as project name.
registry: Docker registry url for images of the project, required.  `,
		Short:   `Create a new project.`,
		Example: `  xiaomei new accounts registry.abc.com/go -d accounts.abc.com`,
		// DisableFlagsInUseLine: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf(`exactly 2 arguments required.`)
			}
			return New(typ, args[0], args[1], domain, force)
		},
	}
	cmd.Flags().StringVarP(&typ, `type`, `t`, `app`, `project type.
 app: only service that provides Golang API.
 web: only service that provides frontend UI.
logc: only service that collect logs to ElasticSearch.
`)
	cmd.Flags().StringVarP(&domain, `domain`, `d`, ``, `domain for the project, only for app and web project, ignored by logc project.
	Used for config.yml, access.conf.tmpl, readme.md, .gitlab-ci.yml.`)
	cmd.Flags().BoolVarP(&force, `force`, `f`, false, `force overwrite existing files.`)
	return cmd
}

func New(typ, projectPath, registry, domain string, force bool) error {
	config, err := getConfig(typ, projectPath, registry, domain)
	if err != nil {
		return err
	}
	templateDir, err := getTemplateDir(typ)
	if err != nil {
		return err
	}
	projectDir := projectPath
	if typ == "app" {
		projectDir = filepath.Base(projectPath)
	}

	return walk(templateDir, projectDir, config, force)
}

func getTemplateDir(typ string) (string, error) {
	modulePath := filepath.Join(fs.GoModPath(), `github.com`, `lovego`, `xiaomei`)
	dirs, err := filepath.Glob(modulePath + `@*`)
	if err != nil {
		return ``, err
	}
	if len(dirs) == 0 {
		if err := utils.GoGetMod(`github.com/lovego/xiaomei`); err != nil {
			return ``, err
		}
		if dirs, err = filepath.Glob(modulePath + `@*`); err != nil {
			return ``, err
		}
		if len(dirs) == 0 {
			return ``, fmt.Errorf("no xiaomei module found at %s.", modulePath)
		}
	}
	sort.Strings(dirs)
	latest := dirs[len(dirs)-1]

	return filepath.Join(latest, `new`, `_templates`, typ), nil
}
