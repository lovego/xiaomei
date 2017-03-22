package app

import (
	"github.com/lovego/xiaomei/xiaomei/images/app/db"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `app`,
		Short: `the app server.`,
	}
	cmd.AddCommand(DepsCmd(), copy2vendorCmd())
	cmd.AddCommand(db.Cmds()...)
	return cmd
}
