package main

import (
	"os"
	// "strings"

	"github.com/bughou-go/xiaomei/xiaomei/images/app"
	"github.com/bughou-go/xiaomei/xiaomei/images/web"
	"github.com/bughou-go/xiaomei/xiaomei/new"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/bughou-go/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	appCmd := app.Cmd()
	webCmd := web.Cmd()
	appCmd.AddCommand(commonCmds(`app`)...)
	webCmd.AddCommand(commonCmds(`web`)...)

	root := rootCmd()
	root.AddCommand(appCmd, webCmd)
	root.AddCommand(commonCmds(``)...)
	root.AddCommand(new.Cmd(), versionCmd())
	root.Execute()
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	/*
		if envs := release.Envs(); len(envs) > 0 {
			cmd.Use += ` [` + strings.Join(envs, `|`) + `]`
		}
	*/
	if release.Arg1IsEnv() {
		cmd.SetArgs(os.Args[2:])
	}
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `show xiaomei version.`,
		RunE: z.NoArgCall(func() error {
			println(`xiaomei version 17.3.8`)
			return nil
		}),
	}
}
