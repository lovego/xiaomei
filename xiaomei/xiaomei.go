package main

import (
	"fmt"
	"os"

	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/new"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/lovego/xiaomei/xiaomei/workspace_godoc"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}
	root.AddCommand(serviceCmds()...)
	root.AddCommand(cluster.Cmds()...)
	root.AddCommand(images.Cmds(``)...)
	root.AddCommand(deploy.Cmds(``)...)
	root.AddCommand(access.Cmds(``)...)
	root.AddCommand(
		new.Cmd(), workspace_godoc.Cmd(),
		yamlCmd(), autoCompleteCmd(root), versionCmd(),
	)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `show xiaomei version.`,
		RunE: release.NoArgCall(func() error {
			fmt.Println(`xiaomei version 18.4.23`)
			return nil
		}),
	}
}
