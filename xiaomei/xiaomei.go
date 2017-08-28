package main

import (
	"os"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/new"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	root := rootCmd()
	root.AddCommand(serviceCmds()...)
	root.AddCommand(images.Cmds(``)...)
	root.AddCommand(deploy.Cmds(``)...)
	root.AddCommand(access.Cmds(``)...)
	root.AddCommand(cluster.LsCmd(), new.Cmd(), yamlCmd(), versionCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:          `xiaomei`,
		Short:        `be small and beautiful.`,
		SilenceUsage: true,
	}
	theCmd.AddCommand(shellCmd())
	return theCmd
}

func shellCmd() *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   `shell [<env>]`,
		Short: `enter a node shell.`,
		RunE: release.EnvCall(func(env string) error {
			_, err := cluster.Get(env).Run(filter, cmd.O{}, ``)
			return err
		}),
		SilenceUsage: true,
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	return theCmd
}
