package main

import (
	"fmt"
	"os"

	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/dbs"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/new"
	"github.com/lovego/xiaomei/xiaomei/oam"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/lovego/xiaomei/xiaomei/spec"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:          `xiaomei`,
		Short:        `be small and beautiful.`,
		SilenceUsage: true,
	}
	root.AddCommand(serviceCmds()...)
	root.AddCommand(deploy.Cmds(``)...)
	root.AddCommand(access.Cmd(``))
	root.AddCommand(oam.Cmds(``)...)
	root.AddCommand(dbs.Cmds()...)
	root.AddCommand(images.Cmds(``)...)
	root.AddCommand(spec.Cmd(), new.Cmd(), yamlCmd(), autoCompleteCmd(root), versionCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `show xiaomei version.`,
		RunE: release.NoArgCall(func() error {
			fmt.Println(`xiaomei version 18.10.8`)
			return nil
		}),
	}
}
