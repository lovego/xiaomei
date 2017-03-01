package project

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `project`,
		Short: `the whole project.`,
	}
	cmd.AddCommand(
		NewCmd(),
	)
	return cmd
}
