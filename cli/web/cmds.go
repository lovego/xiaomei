package web

import (
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `web`,
		Short: `the webserver.`,
	}
	cmd.AddCommand(stack.Cmds(`web`)...)
	return cmd
}
