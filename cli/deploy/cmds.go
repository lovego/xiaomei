package deploy

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var filter *string
	cmd := cobra.Command{
		Use:   `deploy`,
		Short: `[deploy] deploy project.`,
		RunE: func(c *cobra.Command, args []string) error {
			return Deploy(c.Flags().Arg(0), *filter)
		},
	}
	filter = cmd.Flags().StringP(`server`, `s`, ``, `match servers by Addr or Tasks.`)

	cmd.AddCommand(&cobra.Command{
		Use:   `clear-tags`,
		Short: `clear deploy tags.`,
		RunE: func(c *cobra.Command, args []string) error {
			return ClearTags()
		},
	})
	return &cmd
}
