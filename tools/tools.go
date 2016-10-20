package main

import (
	"reflect"

	"github.com/bughou-go/xiaomei/tools/assets"
	"github.com/bughou-go/xiaomei/tools/deploy"
	"github.com/bughou-go/xiaomei/tools/setup"
	// "github.com/bughou-go/xiaomei/tools/tasks"
	"github.com/bughou-go/xiaomei/tools/tools"
)

type cmdFunc struct {
	Func         interface{}
	ArgsRequired int
}

var cmdMaps = map[string]cmdFunc{
	`deploy`:     cmdFunc{Func: deploy.Deploy},
	`restart`:    cmdFunc{Func: deploy.Restart},
	`status`:     cmdFunc{Func: deploy.Status},
	`shell`:      cmdFunc{Func: deploy.Shell},
	`run`:        cmdFunc{Func: deploy.Run},
	`update-doc`: cmdFunc{Func: deploy.UpdateDoc, ArgsRequired: 2},

	`setup`:           cmdFunc{Func: setup.Setup},
	`setup-appserver`: cmdFunc{Func: setup.SetupAppServer},
	`setup-nginx`:     cmdFunc{Func: setup.SetupNginx},
	// `setup-mysql`:     cmdFunc{Func: setup.SetupMysql},
	`setup-cron`:  cmdFunc{Func: setup.SetupCron},
	`setup-hosts`: cmdFunc{Func: setup.SetupHosts},

	`mysql`:     cmdFunc{Func: mysql},
	`mysqldump`: cmdFunc{Func: mysqldump},

	`assets`: cmdFunc{Func: assets.Assets},
}

func main() {
	params := tools.Flags()

	// 参数检查
	if len(params) == 0 {
		tools.PrintUsage()
	}
	if _, ok := cmdMaps[params[0]]; !ok {
		tools.PrintUsage()
	}
	run(params[0], params[1:])
}

// 执行命令
func run(name string, argsStr []string) {
	cmd := cmdMaps[name]
	args := tools.Args(name, reflect.TypeOf(cmd.Func), cmd.ArgsRequired, argsStr)
	reflect.ValueOf(cmd.Func).Call(args)
}
