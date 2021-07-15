package main

import (
	"errors"
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

const moduleVersion = `v1.0.6`
const fullVersion = moduleVersion + ` 20210713`

func main() {
	color.NoColor = false
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `Be small and beautiful.`,
	}
	root.PersistentFlags().SortFlags = false

	root.AddCommand(new.Cmd(moduleVersion))
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
			fmt.Println(`xiaomei ` + fullVersion)
			return nil
		}),
	}
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     `update [version]`,
		Short:   `Update to lastest version.`,
		Example: `xiaomei update v1.0.1`,
		RunE: func(c *cobra.Command, args []string) error {
			target := `github.com/lovego/xiaomei`
			switch len(args) {
			case 0:
			case 1:
				target += `@` + args[0]
			default:
				return errors.New(`more than one arguments given.`)
			}

			fmt.Println(`current version ` + fullVersion)
			if err := release.GoGetByProxy(`-u`, target); err != nil {
				return err
			}
			_, err := cmd.Run(cmd.O{}, `xiaomei`, `version`)
			return err
		},
	}

	return cmd
}
