package setup

import (
	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/cli/setup/nginx"
	"github.com/bughou-go/xiaomei/config"
)

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
		case `nginx`:
			nginx.Setup()
		case `appserver`:
			appserver.Setup()
		default:
			panic(`unknow task: ` + task)
		}
	}
}
