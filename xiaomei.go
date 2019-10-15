package main

import (
	"fmt"
	"os"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/access"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/new"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services"
	"github.com/spf13/cobra"
)

const version = `19.10.15`

func main() {
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	root.AddCommand(new.Cmd())
	root.AddCommand(access.Cmd())
	root.AddCommand(services.Cmds()...)
	root.AddCommand(misc.Cmds(root)...)
	root.AddCommand(versionCmd(), updateCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `Show xiaomei version.`,
		RunE: release.NoArgCall(func() error {
			fmt.Println(`xiaomei version ` + version)
			return nil
		}),
	}
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `update`,
		Short: `Update to lastest version.`,
		RunE: release.NoArgCall(func() error {
			fmt.Println(`current version ` + version)
			if _, err := cmd.Run(cmd.O{},
				`go`, `get`, `-u`, `-v`, `github.com/lovego/xiaomei`,
			); err != nil {
				return err
			}
			_, err := cmd.Run(cmd.O{}, `xiaomei`, `version`)
			return err
		}),
	}
}
