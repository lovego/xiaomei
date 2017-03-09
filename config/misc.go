package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/bughou-go/xiaomei/utils"
)

const ISO8601 = `2006-01-02T15:04:05Z0700`

func Log(msg string) {
	println(time.Now().Format(ISO8601), msg)
}

func Logf(w io.Writer, msg string) {
	w.Write([]byte(time.Now().Format(ISO8601) + ` ` + msg + "\n"))
}

func Debug(name string) bool {
	matched, _ := regexp.MatchString(`\b`+name+`\b`, os.Getenv(`debug`))
	return matched
}

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("PANIC: %s\n%s", err, utils.Stack(4))
			Alarm(`Protect错误`, errMsg)
			log.Printf(errMsg)
		}
	}()
	fn()
}
