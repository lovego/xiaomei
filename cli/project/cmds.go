package project

import (
	new "github.com/bughou-go/xiaomei/cli/project/new"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `pj`,
		Short: `the whole project.`,
	}
	cmd.AddCommand(stack.BDPcmds(``)...)
	cmd.AddCommand(new.Cmd())
	return cmd
}
