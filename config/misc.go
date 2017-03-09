package config

import (
	"fmt"

	"github.com/bughou-go/xiaomei/utils"
)

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
