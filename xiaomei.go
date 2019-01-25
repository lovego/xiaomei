package main

import (
	"fmt"
	"os"

	"github.com/lovego/xiaomei/access"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/new"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:          `xiaomei`,
		Short:        `be small and beautiful.`,
		SilenceUsage: true,
	}

	root.AddCommand(new.Cmd())
	root.AddCommand(services.Cmds()...)
	root.AddCommand(access.Cmd())
	root.AddCommand(misc.Cmds(root)...)
	root.AddCommand(versionCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `show xiaomei version.`,
		RunE: release.NoArgCall(func() error {
			fmt.Println(`xiaomei version 19.01.24`)
			return nil
		}),
	}
}
