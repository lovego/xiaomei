package logc

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `logc`,
		Short: `the log collector client.`,
	}
	return cmd
}
