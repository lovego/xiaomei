package main

import (
	"errors"
	"os"

	"github.com/lovego/xiaomei/utils/cmd"
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
	root.AddCommand(cluster.LsCmd(), new.Cmd(), yamlCmd(), versionCmd())
	root.Execute()
}

func rootCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `xiaomei [qa|production|...]`,
		Short: `be small and beautiful.`,
		RunE: func(thisCmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New(`redundant args.`)
			}
			if release.Arg1IsEnv() {
				_, err := cluster.Run(cmd.O{}, ``)
				return err
			} else {
				return thisCmd.Help()
			}
		},
	}
	if release.Arg1IsEnv() {
		theCmd.SetArgs(os.Args[2:])
	}
	return theCmd
}
