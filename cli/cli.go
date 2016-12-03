package cli

import (
	"reflect"

	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/cli/db"
	"github.com/bughou-go/xiaomei/cli/deploy"
	"github.com/bughou-go/xiaomei/cli/develop"
	"github.com/bughou-go/xiaomei/cli/oam"
	"github.com/bughou-go/xiaomei/cli/setup"
)

type CmdFunc struct {
	Func         interface{}
	ArgsRequired int
}

var CmdMaps = map[string]CmdFunc{
	`new`:    CmdFunc{Func: develop.New},
	`run`:    CmdFunc{Func: develop.Run},
	`build`:  CmdFunc{Func: develop.Build},
	`spec`:   CmdFunc{Func: develop.Spec},
	`assets`: CmdFunc{Func: develop.Assets},

	`deploy`:     CmdFunc{Func: deploy.Deploy},
	`restart`:    CmdFunc{Func: deploy.Restart},
	`status`:     CmdFunc{Func: deploy.Status},
	`shell`:      CmdFunc{Func: deploy.Shell},
	`exec`:       CmdFunc{Func: deploy.Exec},
	`update-doc`: CmdFunc{Func: deploy.UpdateDoc, ArgsRequired: 2},

	`setup`:           CmdFunc{Func: setup.Setup},
	`setup-appserver`: CmdFunc{Func: setup.SetupAppServer},
	`setup-nginx`:     CmdFunc{Func: setup.SetupNginx},
	`setup-cron`:      CmdFunc{Func: setup.SetupCron},
	`setup-hosts`:     CmdFunc{Func: setup.SetupHosts},
	//`setup-mysql`:     CmdFunc{Func: setup.SetupMysql},

	`mysql`:     CmdFunc{Func: db.Mysql},
	`mysqldump`: CmdFunc{Func: db.Mysqldump},
}

func Run() {
	params := cli.Flags()

	// 参数检查
	if len(params) == 0 {
		cli.PrintUsage()
	}
	if _, ok := CmdMaps[params[0]]; !ok {
		cli.PrintUsage()
	}
	run(params[0], params[1:])
}

// 执行命令
func run(name string, argsStr []string) {
	cmd := CmdMaps[name]
	args := cli.Args(name, reflect.TypeOf(cmd.Func), cmd.ArgsRequired, argsStr)
	reflect.ValueOf(cmd.Func).Call(args)
}
