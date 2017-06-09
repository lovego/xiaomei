package main

import (
	"os"

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

func manageCmds() []*cobra.Command {
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

	logcCmd := logc.Cmd()
	logcCmd.AddCommand(deploy.Cmds(`logc`)...)

	return []*cobra.Command{appCmd, tasksCmd, webCmd, godocCmd, logcCmd, cluster.Cmd()}
}
