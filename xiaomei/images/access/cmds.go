package access

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access`,
		Short: `the access server.`,
	}
	return cmd
}
