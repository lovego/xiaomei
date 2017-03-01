package main

import (
	"github.com/bughou-go/xiaomei/cli/app"
	"github.com/bughou-go/xiaomei/cli/project"
	// "github.com/bughou-go/xiaomei/cli/db"

	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	root.AddCommand(
		app.Cmd(), project.Cmd(),
	)
	/*
		root.AddCommand(db.Cmds()...)
	*/

	root.Execute()
}
