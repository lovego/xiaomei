package app

import (
	"github.com/bughou-go/xiaomei/cli/app/deps"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `app`,
		Short: `the appserver.`,
	}
	cmd.AddCommand(runCmd())
	cmd.AddCommand(stack.BDPcmds(`app`)...)
	cmd.AddCommand(deps.Cmd(), SpecCmd())
	return cmd
}
