package deploy

import (
	"github.com/spf13/cobra"
)

func Cmd(serverFilter *string) *cobra.Command {
	cmd := cobra.Command{
		Use:   `deploy`,
		Short: `deploy project to a environment.`,
		RunE: func(c *cobra.Command, args []string) error {
			return Deploy(c.Flags().Arg(0), *serverFilter)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   `clear-tags`,
		Short: `clear deploy tags.`,
		RunE: func(c *cobra.Command, args []string) error {
			return ClearTags()
		},
	})
	return &cmd
}
