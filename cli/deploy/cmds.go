package deploy

import (
	"github.com/spf13/cobra"
)

func Cmds(serverFilter *string) []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `deploy`,
			Short: `deploy project to a environment.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Deploy(c.Flags().Arg(0), *serverFilter)
			},
		},
	}
}
