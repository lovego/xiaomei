package db

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use: `mysql`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mysql(name)
			},
		},
		{
			Use: `mysqldump`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mysqldump(name)
			},
		},
	}
}
