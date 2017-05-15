package godoc

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `godoc`,
		Short: `the godoc server.`,
	}
	return cmd
}
