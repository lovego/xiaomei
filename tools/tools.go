package main

import (
	"fmt"
	"reflect"
	"github.com/bughou-go/xiaomei/tools/assets"
	"github.com/bughou-go/xiaomei/tools/deploy"
	"github.com/bughou-go/xiaomei/tools/setup"
	"github.com/bughou-go/xiaomei/tools/tasks"
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
	`setup-mysql`:     cmdFunc{Func: setup.SetupMysql},
	`setup-cron`:      cmdFunc{Func: setup.SetupCron},
	`setup-hosts`:     cmdFunc{Func: setup.SetupHosts},

	`mysql`:     cmdFunc{Func: mysql},
	`mysqldump`: cmdFunc{Func: mysqldump},

	`sync-orgs`: cmdFunc{Func: tasks.SyncOrgs},

	`assets`: cmdFunc{Func: assets.Assets},

	`plan-users`: cmdFunc{Func: tasks.CountPlanUsers},

	`baidu-weather`: cmdFunc{Func: tasks.GetUpdateWeather},
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
	args := getArgs(name, reflect.TypeOf(cmd.Func), cmd.ArgsRequired, argsStr)
	reflect.ValueOf(cmd.Func).Call(args)
}

func getArgs(name string, function reflect.Type, required int, args []string) (
	funcArgs []reflect.Value,
) {
	optional, variadic := cmdArgsNum(name, required, function)
	length := len(args)
	if length < required {
		panic(fmt.Sprintf("命令%s: 至少需要%d个参数, 但得到了%v个。", name, required, length))
	}
	if !variadic && length > (required+optional) {
		panic(fmt.Sprintf("命令%s: 最多需要%d个参数, 但得到了%v个。", name, required+optional, length))
	}
	funcArgs = appendRequiredArgs(required, args, funcArgs)
	funcArgs = appendOptionalArgs(required, optional, args, funcArgs)
	funcArgs = appendVariadicArgs(required+optional, variadic, args, funcArgs)
	return
}

func appendRequiredArgs(required int, args []string, funcArgs []reflect.Value) []reflect.Value {
	for i := 0; i < required; i++ {
		funcArgs = append(funcArgs, reflect.ValueOf(args[i]))
	}
	return funcArgs
}

func appendOptionalArgs(required, optional int, args []string, funcArgs []reflect.Value,
) []reflect.Value {
	length := len(args)
	for i := 0; i < optional; i++ {
		if required+i < length {
			funcArgs = append(funcArgs, reflect.ValueOf(args[required+i]))
		} else {
			funcArgs = append(funcArgs, reflect.ValueOf(``))
		}
	}
	return funcArgs
}

func appendVariadicArgs(l int, variadic bool, args []string, funcArgs []reflect.Value,
) []reflect.Value {
	if variadic {
		if len(args) > l {
			funcArgs = append(funcArgs, reflect.ValueOf(args[l:]))
		} else {
			funcArgs = append(funcArgs, reflect.ValueOf([]string{}))
		}
	}
	return funcArgs
}

var typOptional = reflect.TypeOf(``)
var typVariadic = reflect.TypeOf([]string{})

// 参数类型顺序： 必选，可选，可变
func cmdArgsNum(name string, required int, function reflect.Type) (optional int, variadic bool) {
	last := function.NumIn() - 1
	i := required
	for ; i <= last && function.In(i) == typOptional; i++ {
	}
	optional = i - required

	switch {
	case i > last:
	case i == last:
		if function.In(i) == typVariadic {
			variadic = true
		} else {
			panic(fmt.Sprintf("命令%s: 第%d个参数类型非法", name, i))
		}
	default:
		panic(fmt.Sprintf("命令%s: 第%d个参数类型非法", name, i))
	}
	return
}
