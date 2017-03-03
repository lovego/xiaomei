package main

import (
	"os"
	"strings"

	"github.com/bughou-go/xiaomei/cli/app"
	"github.com/bughou-go/xiaomei/cli/project"
	"github.com/bughou-go/xiaomei/config"
	// "github.com/bughou-go/xiaomei/cli/db"

	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	root := rootCmd()

	root.AddCommand(
		app.Cmd(), project.Cmd(),
	)
	/*
		root.AddCommand(db.Cmds()...)
	*/

	root.Execute()
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}
	if envs := config.Envs(); len(envs) > 0 {
		cmd.Use += `[` + strings.Join(envs, `|`) + `]`
	}
	if config.Arg1IsEnv() {
		cmd.SetArgs(os.Args[2:])
	}
	return cmd
}
