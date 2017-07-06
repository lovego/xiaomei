package main

import (
	"errors"
	"os"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/godoc"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/tasks"
	"github.com/lovego/xiaomei/xiaomei/images/web"
	"github.com/lovego/xiaomei/xiaomei/new"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	root := rootCmd()
	root.AddCommand(manageCmds()...)
	root.AddCommand(images.Cmds(``)...)
	root.AddCommand(deploy.Cmds(``)...)
	root.AddCommand(cluster.LsCmd(), new.Cmd(), yamlCmd(), versionCmd())
	root.Execute()
}

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   `xiaomei [qa|production|...]`,
		Short: `be small and beautiful.`,
		RunE: func(rootCmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New(`redundant args.`)
			}
			if release.Arg1IsEnv() {
				_, err := cluster.Run(cmd.O{}, ``)
				return err
			} else {
				return rootCmd.Help()
			}
		},
	}
	if release.Arg1IsEnv() {
		rootCmd.SetArgs(os.Args[2:])
	}
	return rootCmd
}

func manageCmds() []*cobra.Command {
	appCmd := app.Cmd()
	appCmd.AddCommand(images.Cmds(`app`)...)
	appCmd.AddCommand(deploy.Cmds(`app`)...)

	tasksCmd := tasks.Cmd()
	tasksCmd.AddCommand(images.Cmds(`tasks`)...)
	tasksCmd.AddCommand(deploy.Cmds(`tasks`)...)

	webCmd := web.Cmd()
	webCmd.AddCommand(images.Cmds(`web`)...)
	webCmd.AddCommand(deploy.Cmds(`web`)...)

	logcCmd := logc.Cmd()
	logcCmd.AddCommand(images.Cmds(`logc`)...)
	logcCmd.AddCommand(deploy.Cmds(`logc`)...)

	godocCmd := godoc.Cmd()
	godocCmd.AddCommand(images.Cmds(`godoc`)...)
	godocCmd.AddCommand(deploy.Cmds(`godoc`)...)

	return []*cobra.Command{appCmd, tasksCmd, webCmd, logcCmd, godocCmd}
}
