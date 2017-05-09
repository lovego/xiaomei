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

func Alarm(title, body, mergeKey string) {
	keepers := Keepers()
	if len(keepers) == 0 {
		return
	}
	title = DeployName() + ` ` + title
	err := Mailer.Alarm(&email.Email{
		To:      keepers,
		Subject: title,
		Text:    []byte(body),
	}, mergeKey)
	if err != nil {
		utils.Log(`send alarm mail failed: ` + err.Error())
	}
}

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			errStack := fmt.Sprintf("PANIC: %s\n%s", err, utils.Stack(4))
			errMsg := time.Now().Format(utils.ISO8601) + ` ` + errStack
			go Alarm(`Protect错误`, errMsg, errStack)
			utils.Log(errMsg)
		}
	}()
	fn()
}
