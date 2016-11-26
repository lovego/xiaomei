package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/mailer"
)

var Mailer *mailer.Mailer

func setupMailer() {
	m := Data.Mailer
	if m.Host == `` || m.Port == `` || m.Sender == `` {
		return
	}
	Mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
}

func AlarmMail(title, body string) {
	title = Data.DeployName + ` ` + title
	Mailer.Send(&mailer.Message{Receivers: Data.AlarmReceivers, Title: title, Body: body})
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
			AlarmMail(`Protect错误`, errMsg)
			log.Printf(errMsg)
		}
	}()
	fn()
}
