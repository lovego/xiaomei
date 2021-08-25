package new

import (
	"fmt"
	"path/filepath"

	"github.com/lovego/cmd"
	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd(moduleVersion string) *cobra.Command {
	var typ, domain string
	var force bool
	cmd := &cobra.Command{
		Use: `new <project dir> <registry> [flags]

project dir: Go module path or directory for the project, required. The last element of project path is used as project name.
registry: Docker registry url for images of the project, required.  `,
		Short:   `Create a new project.`,
		Example: `  xiaomei new accounts registry.abc.com/backend -d accounts.abc.com`,
		// DisableFlagsInUseLine: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf(`exactly 2 arguments required.`)
			}
			return New(moduleVersion, typ, args[0], args[1], domain, force)
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

func New(moduleVersion, typ, projectDir, registry, domain string, force bool) error {
	config, err := getConfig(typ, projectDir, registry, domain)
	if err != nil {
		return err
	}
	templateDir, err := getTemplateDir(moduleVersion, typ)
	if err != nil {
		return err
	}

	if err := walk(templateDir, projectDir, config, force); err != nil {
		return err
	}
	if typ == `app` {
		if err := initProject(projectDir); err != nil {
			return err
		}
	}
	return nil
}

func getTemplateDir(moduleVersion, typ string) (string, error) {
	moduleDir := filepath.Join(fs.GoModPath(), `github.com`, `lovego`, `xiaomei`+`@`+moduleVersion)
	if !fs.IsDir(moduleDir) {
		module := `github.com/lovego/xiaomei` + `@` + moduleVersion
		fmt.Println(module, `download...`)
		if err := release.GoGetByProxy(module); err != nil {
			return ``, err
		}
		if !fs.IsDir(moduleDir) {
			return ``, fmt.Errorf("no xiaomei module found at %s.", moduleDir)
		}
	}
	return filepath.Join(moduleDir, `new`, `_templates`, typ), nil
}

func initProject(projectDir string) error {
	o := cmd.O{Dir: filepath.Join(projectDir, "src")}
	if _, err := cmd.Run(o, release.GoCmd(), `mod`, `init`, projectDir); err != nil {
		return err
	}
	if _, err := cmd.Run(o, release.GoCmd(), `mod`, `vendor`); err != nil {
		return err
	}
	return nil
}
