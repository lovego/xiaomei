package main

import (
	"github.com/bughou-go/xiaomei/cli/app"
	"github.com/bughou-go/xiaomei/cli/project"
	// "github.com/bughou-go/xiaomei/cli/db"
	// "github.com/bughou-go/xiaomei/cli/oam"

	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	root.AddCommand(
		app.Cmd(), project.NewCmd(),
	)
	/*
		root.AddCommand(db.Cmds()...)
		root.AddCommand(oam.Cmds()...)
	*/

	root.Execute()
}
