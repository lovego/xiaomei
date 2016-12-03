package cli

import (
	"fmt"
	"reflect"
)

func Args(cmd string, function reflect.Type, required int, args []string) (
	funcArgs []reflect.Value,
) {
	optional, variadic := cmdArgsNum(cmd, required, function)
	length := len(args)
	if length < required {
		panic(fmt.Sprintf("命令%s: 至少需要%d个参数, 但得到了%v个。", cmd, required, length))
	}
	if !variadic && length > (required+optional) {
		panic(fmt.Sprintf("命令%s: 最多需要%d个参数, 但得到了%v个。", cmd, required+optional, length))
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
func cmdArgsNum(cmd string, required int, function reflect.Type) (optional int, variadic bool) {
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
			panic(fmt.Sprintf("命令%s: 第%d个参数类型非法", cmd, i))
		}
	default:
		panic(fmt.Sprintf("命令%s: 第%d个参数类型非法", cmd, i))
	}
	return
}
