package setup

import (
	"github.com/bughou-go/xiaomei/config"
	"net"
	"regexp"
	"strings"
)

func Setup(tasks string) {
	if tasks == `` {
		tasks = currentServerTasks()
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

func currentServerTasks() (tasks string) {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range config.Data().DeployServers {
	loop:
		for _, ifcAddr := range ifcAddrs {
			if strings.HasPrefix(ifcAddr.String(), server.Addr+`/`) {
				tasks += ` ` + server.Tasks
				break loop
			}
		}
	}
	return
}
