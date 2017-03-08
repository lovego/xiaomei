package app

import (
	"github.com/bughou-go/xiaomei/xiaomei/images/app/deps"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `app`,
		Short: `the app server.`,
	}
	cmd.AddCommand(deps.Cmd(), SpecCmd())
	return cmd
}
