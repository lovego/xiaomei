package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/lovego/xiaomei/utils"
)

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

func Alarm(title, body string) {
	keepers := Keepers()
	if len(keepers) == 0 {
		return
	}
	title = DeployName() + ` ` + title
	err := Mailer.Send(&email.Email{
		To:      keepers,
		Subject: title,
		Text:    []byte(body),
	}, time.Minute)
	utils.Log(err.Error())
}

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			errMsg := fmt.Sprintf("PANIC: %s\n%s", err, utils.Stack(4))
			Alarm(`Protect错误`, errMsg)
			utils.Log(errMsg)
		}
	}()
	fn()
}
