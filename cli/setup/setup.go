package setup

import (
	"regexp"

	"github.com/bughou-go/xiaomei/config"
)

func Setup(tasks string) {
	if tasks == `` {
		tasks = config.Servers.CurrentTasks()
	}
	for _, task := range regexp.MustCompile(`[\w-]+`).FindAllString(tasks, -1) {
		switch task {
		case `setup-hosts`:
			SetupHosts()
		case `setup-mysql`:
			// SetupMysql()
		case `setup-appserver`:
			SetupAppServer()
		case `setup-cron`:
			SetupCron()
		case `setup-nginx`:
			SetupNginx()
		default:
			panic(`unknow task: ` + task)
		}
	}
}
