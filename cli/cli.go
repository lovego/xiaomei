package cli

import (
	"errors"

	"github.com/bughou-go/xiaomei/cli/db"
	"github.com/bughou-go/xiaomei/cli/deploy"
	"github.com/bughou-go/xiaomei/cli/develop"
	"github.com/bughou-go/xiaomei/cli/oam"
	"github.com/bughou-go/xiaomei/cli/setup"

	"github.com/spf13/cobra"
)

func Run() {
	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}
	flags := root.PersistentFlags()
	s := flags.StringP(
		`server`, `s`, ``, `match servers by Addr or Tasks.`,
	)
	root.AddCommand(develop.Cmds()...)
	root.AddCommand(deploy.Cmds(s)...)
	root.AddCommand(setup.Cmds()...)
	root.AddCommand(oam.Cmds(s)...)
	root.AddCommand(db.Cmds()...)

	root.Execute()
}

var argsCountError = errors.New(`args count doesn't match.`)

func checkArgsNum(c *cobra.Command, args []string) {
	// c.Use
}
