package config

import (
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/alarm"
)

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

var alarmEngine = alarm.NewEngine(DeployName(), alarm.MailSender{
	Receivers: Keepers(),
	Mailer:    Mailer(),
}, 0, time.Second, 10*time.Second, nil)

func AlarmEngine() *alarm.Engine {
	return alarmEngine
}

func Protect(fn func()) {
	defer alarmEngine.Recover()
	fn()
}
