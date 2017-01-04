package cli

import (
	"github.com/bughou-go/xiaomei/cli/db"
	"github.com/bughou-go/xiaomei/cli/deploy"
	"github.com/bughou-go/xiaomei/cli/oam"
	"github.com/bughou-go/xiaomei/cli/setup"
	"github.com/bughou-go/xiaomei/config"

	"github.com/spf13/cobra"
)

func Run() {
	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:   config.App.Name(),
		Short: `setup the specified tasks.`,
	}
	root.AddCommand(db.Cmds()...)
	root.AddCommand(oam.Cmds()...)
	root.AddCommand(deploy.Cmd(), setup.Cmd())

	root.Execute()
}
