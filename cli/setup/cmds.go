package setup

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `setup [<tasks>]`,
			Short: `setup all tasks.`,
			Run: func(c *cobra.Command, args []string) {
				var tasks string
				if len(args) > 0 {
					tasks = args[0]
				}
				Setup(tasks)
			},
		},
		{
			Use:   `setup-appserver`,
			Short: `setup appserver.`,
			Run: func(c *cobra.Command, args []string) {
				SetupAppServer()
			},
		},
		{
			Use: `setup-nginx`,
			Run: func(c *cobra.Command, args []string) {
				SetupNginx()
			},
		},
		{
			Use:   `setup-docserver`,
			Short: `start doc-server.`,
			Run: func(c *cobra.Command, args []string) {
				DocServer()
			},
		},
		{
			Use: `setup-cron`,
			Run: func(c *cobra.Command, args []string) {
				SetupCron()
			},
		},
		{
			Use: `setup-hosts`,
			Run: func(c *cobra.Command, args []string) {
				SetupHosts()
			},
		},
	}
}
