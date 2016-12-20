package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/bughou-go/xiaomei/utils"
)

func Debug(name string) bool {
	matched, _ := regexp.MatchString(`\b`+name+`\b`, os.Getenv(`debug`))
	return matched
}

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("PANIC: %s\n%s", err, utils.Stack(4))
			App.Alarm(`Protect错误`, errMsg)
			log.Printf(errMsg)
		}
	}()
	fn()
}
