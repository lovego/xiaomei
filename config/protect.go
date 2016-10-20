package config

import (
	"fmt"
	"log"
	"github.com/bughou-go/xm"
)

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			err_message := fmt.Sprintf("PANIC: %s\n%s", err, xm.Stack(4))
			AlarmMail(`Protect错误`, err_message)
			log.Printf(err_message)
		}
	}()
	fn()
}
