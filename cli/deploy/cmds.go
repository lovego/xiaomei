package deploy

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `deploy`,
			Short: `deploy project to a environment.`,
			RunE: func(c *cobra.Command, args []string) error {
				var tasks string
				if len(args) > 0 {
					tasks = args[0]
				}
				return Deploy(tasks)
			},
		},
	}
}
