package config

import (
	"fmt"
	"os"

	"github.com/bughou-go/xiaomei/utils"
)

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("PANIC: %s\n%s", err, utils.Stack(4))
			Alarm(`Protect错误`, errMsg)
			println(errMsg)
		}
	}()
	fn()
}
