package setup

import (
	"errors"

	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/cli/setup/godoc"
	"github.com/bughou-go/xiaomei/cli/setup/nginx"
	"github.com/bughou-go/xiaomei/config"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `setup [nginx|appserver|hosts|mysql|cron|godoc] ...`,
			Short: `setup nginx, appserver, hosts, mysql, cron and godoc.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Setup(args)
			},
		},
		{
			Use:    `launch`,
			Short:  `launch appserver.`,
			Hidden: true,
			Run: func(c *cobra.Command, args []string) {
				appserver.Launch()
			},
		},
	}
}

func Setup(tasks []string) error {
	if len(tasks) == 0 {
		tasks = config.Servers.CurrentTasks()
	}
	for _, task := range tasks {
		switch task {
		case `hosts`:
			SetupHosts()
		case `mysql`:
			SetupMysql()
		case `cron`:
			SetupCron()
		case `godoc`:
			godoc.InDeploy()
		case `nginx`:
			nginx.Setup()
			godoc.SetupNginxInDeploy()
		case `appserver`:
			if err := appserver.Restart(true); err != nil {
				return err
			}
			godoc.SetupUpstartInDeploy()
		default:
			return errors.New(`unknown task: ` + task)
		}
	}
	return nil
}
