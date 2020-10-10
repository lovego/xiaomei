package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/access"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/new"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services"
	"github.com/spf13/cobra"
)

const version = `20.10.10`

func main() {
	color.NoColor = false
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `Be small and beautiful.`,
	}
	root.PersistentFlags().SortFlags = false

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
	var deep bool
	cmd := &cobra.Command{
		Use:   `update`,
		Short: `Update to lastest version.`,
		RunE: release.NoArgCall(func() error {
			fmt.Println(`current version ` + version)
			if deep {
				if _, err := cmd.Run(cmd.O{},
					`go`, `get`, `-u`, `-v`, `github.com/lovego/xiaomei`,
				); err != nil {
					return err
				}
			} else if _, err := cmd.Run(cmd.O{}, `sh`, `-c`, `
				cd $(go env GOPATH)/src/github.com/lovego/xiaomei && git pull && go install -v
			`); err != nil {
				return err
			}
			_, err := cmd.Run(cmd.O{}, `xiaomei`, `version`)
			return err
		}),
	}
	cmd.Flags().BoolVar(&deep, `deep`, false, `Update all dependency packages.`)

	return cmd
}
