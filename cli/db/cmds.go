package db

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `mysql [dbname | default]`,
			Short: `into the specified dbname of mysql command line tool`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mysql(name)
			},
		},
		{
			Use:   `mysqldump [dbname | default]`,
			Short: `dump the specified dbname of mysql`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mysqldump(name)
			},
		},
		{
			Use:   `mongo [dbname | default]`,
			Short: `into the specified dbname of mongo command line tool`,
			Run: func(c *cobra.Command, args []string) {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				Mongo(name)
			},
		},
		{
			Use:   `redis [dbname | default]`,
			Short: `into the specified dbname of redis command line tool`,
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
