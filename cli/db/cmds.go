package db

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `mysql [<dbkey>]`,
			Short: `[db] enter mysql cli`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mysql(name)
			},
		},
		{
			Use:   `mysqldump [<dbkey>]`,
			Short: `[db] dump mysql`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mysqldump(name)
			},
		},
		{
			Use:   `mongo [<dbkey>]`,
			Short: `[db] enter mongo cli`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mongo(name)
			},
		},
		{
			Use:   `redis [<dbkey>]`,
			Short: `[db] enter redis cli`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Redis(name)
			},
		},
	}
}
