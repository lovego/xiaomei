package main

import (
	"os"

	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/access"
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

	appCmd := app.Cmd()
	appCmd.AddCommand(images.Cmds(`app`)...)
	appCmd.AddCommand(deploy.Cmds(`app`)...)

	tasksCmd := tasks.Cmd()
	tasksCmd.AddCommand(images.Cmds(`tasks`)...)
	tasksCmd.AddCommand(deploy.Cmds(`tasks`)...)

	godocCmd := godoc.Cmd()
	godocCmd.AddCommand(images.Cmds(`godoc`)...)
	godocCmd.AddCommand(deploy.Cmds(`godoc`)...)

	webCmd := web.Cmd()
	webCmd.AddCommand(images.Cmds(`web`)...)
	webCmd.AddCommand(deploy.Cmds(`web`)...)

	accessCmd := access.Cmd()
	accessCmd.AddCommand(images.Cmds(`access`)...)
	accessCmd.AddCommand(deploy.Cmds(`access`)...)

	logcCmd := logc.Cmd()
	logcCmd.AddCommand(deploy.Cmds(`logc`)...)

	root := rootCmd()
	root.AddCommand(appCmd, tasksCmd, godocCmd, webCmd, accessCmd, logcCmd, cluster.Cmd())
	root.AddCommand(images.Cmds(``)...)
	root.AddCommand(deploy.Cmds(``)...)
	root.AddCommand(new.Cmd(), yamlCmd(), versionCmd())
	root.Execute()
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	if release.Arg1IsEnv() {
		cmd.SetArgs(os.Args[2:])
	}
	return cmd
}
