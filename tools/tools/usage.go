package tools

import (
	"fmt"
	"os"
)

func PrintUsage() {
	fmt.Printf(`Usage:
  %s [-s server] command ...

  The commands are:
  deploy            部署
  restart           重启应用服务器
  status            查看应用服务器状态
  shell             进入服务器的bash
  run <cmd>         在服务器运行<cmd>

  setup           设置mysql、appserver、cron、nginx
  setup-appserver 设置应用服务器
  setup-nginx     设置nginx
  setup-mysql     在mysql中建库建表，导入初始数据
  setup-cron      设置定时任务(/etc/cron.d)
  setup-hosts     设置hosts文件(/etc/hosts)

  mysql 进入mysql命令客户端
  mysqldump 执行mysqldump命令导出mysql数据

  assets [args...] 添加至assets.json

  sync-orgs 同步采销Tree组织
  plan-users 获取当天plan使用的用户
  baidu-weather 获取并更新当天及后三天的天气数据
`, os.Args[0])
	os.Exit(1)
}
