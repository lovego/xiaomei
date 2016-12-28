package setup

import (
	"github.com/bughou-go/xiaomei/cli/godoc"
	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/cli/setup/nginx"
	"github.com/bughou-go/xiaomei/config"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   `setup [<tasks> ...]`,
		Short: `setup the specified tasks.`,
		Run: func(c *cobra.Command, args []string) {
			Setup(args)
		},
	}
}

func Setup(tasks []string) {
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
			appserver.Setup()
			godoc.SetupUpstartInDeploy()
		case `wait-appserver`:
			appserver.Wait()
		default:
			panic(`unknown task: ` + task)
		}
	}
}
