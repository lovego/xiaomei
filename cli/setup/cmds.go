package setup

import (
	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/cli/setup/nginx"
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
				appserver.Setup()
			},
		},
		{
			Use:    `wait-appserver`,
			Short:  `wait appserver until it's started.`,
			Hidden: true,
			Run: func(c *cobra.Command, args []string) {
				appserver.Wait()
			},
		},
		{
			Use:   `setup-mysql`,
			Short: `setup mysql.`,
			Run: func(c *cobra.Command, args []string) {
				SetupMysql()
			},
		},
		{
			Use:   `setup-nginx`,
			Short: `setup nginx.`,
			Run: func(c *cobra.Command, args []string) {
				nginx.Setup()
			},
		},
		{
			Use:   `setup-docserver`,
			Short: `setup doc-server.`,
			Run: func(c *cobra.Command, args []string) {
				DocServer()
			},
		},
		{
			Use:   `setup-cron`,
			Short: `setup crontab.`,
			Run: func(c *cobra.Command, args []string) {
				SetupCron()
			},
		},
		{
			Use:   `setup-hosts`,
			Short: `setup hosts.`,
			Run: func(c *cobra.Command, args []string) {
				SetupHosts()
			},
		},
	}
}
