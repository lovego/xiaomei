package config

import (
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/alarm"
	"github.com/lovego/xiaomei/utils/logger"
)

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

var theAlarm = alarm.New(DeployName(), alarm.MailSender{
	Receivers: Keepers(),
	Mailer:    Mailer(),
}, 0, time.Second, 10*time.Second)

var theLogger = logger.New(``, os.Stderr, theAlarm)

func Alarm() *alarm.Alarm {
	return theAlarm
}

func Logger() *logger.Logger {
	return theLogger
}

func Protect(fn func()) {
	defer theLogger.Recover()
	fn()
}
