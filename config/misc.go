package config

import (
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/alarm"
)

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

var theAlarm = alarm.New(DeployName(), alarm.MailSender{
	Receivers: Keepers(),
	Mailer:    Mailer(),
}, 0, time.Second, 10*time.Second, nil)

func Alarm() *alarm.Alarm {
	return theAlarm
}

func Protect(fn func()) {
	defer theAlarm.Recover()
	fn()
}
