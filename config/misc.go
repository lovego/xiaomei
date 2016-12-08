package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/mailer"
)

var timeZone *time.Location

func TimeZone() *time.Location {
	if timeZone == nil {
		time.FixedZone(data().TimeZoneName, data().TimeZoneOffset)
	}
	return timeZone
}

var _mailer *mailer.Mailer
var mailerSet bool
var mailerMutex sync.Mutex

func Mailer() *mailer.Mailer {
	mailerMutex.Lock()
	defer mailerMutex.Unlock()
	if !mailerSet {
		m := data().Mailer
		if m.Host != `` || m.Port != `` || m.Sender != `` {
			_mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
		}
		mailerSet = true
	}
	return _mailer
}

func AlarmMail(title, body string) {
	title = DeployName() + ` ` + title
	Mailer().Send(&mailer.Message{Receivers: data().AlarmReceivers, Title: title, Body: body})
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
