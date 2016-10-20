package tasks

import (
	"fmt"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"time"
)

func CountPlanUsers() {
	data := make(map[string]string)
	date := time.Now().In(config.TimeZone)
	cur_addr := config.CurrentAppServer().Addr
	for _, addr := range getServersAddrs() {
		runCmds(data, addr, date, cur_addr == addr)
	}
	config.SendMailWithAttachs([]string{`赵佳<zhaojiaa@wumart.com>`},
		date.Format(`2006-01-02 15:04:05`)+`统计 app和PC计划管理访问清单`, `详细见附件`, data)
}

func getServersAddrs() []string {
	addrs := []string{}
	for _, deploy_server := range config.Data.DeployServers {
		if deploy_server.AppAddr == `` {
			addr := deploy_server.Addr
			addrs = append(addrs, addr)
		}
	}
	return addrs
}

func runCmds(data map[string]string, addr string, date time.Time, is_self bool) {
	cur_date := date.Format(`2006-01-02`)
	last_date := date.AddDate(0, 0, -1).Format(`2006-01-02`)
	names := []string{`plan`, `plan_app`, `vendor_plan`, `kpi`}
	commands := []string{
		`egrep -ah '(%s|%s) .* /plan/index/?[ ?#]' %s/release/log/app.log`,
		`egrep -ah '(%s|%s) .* /api/category_plan/cxdata/' %s/release/log/app.log`,
		`egrep -ah '(%s|%s) .* /vendor/plan/?[ ?#]' %s/release/log/app.log`,
		`egrep -ah '(%s|%s) .* /kpi/?[ ?#]' %s/release/log/app.log`,
	}

	for i, command := range commands {
		var output string
		command = fmt.Sprintf(command, cur_date, last_date, config.Data.DeployPath)
		if is_self {
			output, _ = cmd.Run(cmd.O{Output: true}, ``, command)
		} else {
			output, _ = cmd.Run(cmd.O{Output: true}, `ssh`, config.Data.DeployUser+`@`+addr, command)
		}
		name := fmt.Sprintf("%s_%s.txt", names[i], cur_date)
		data[name] += fmt.Sprintf("\n%s", output)
	}
}
