package config

import (
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/alarm"
)

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

var alarmEngine = alarm.NewEngine(alarm.MailSender{
	Receivers: Keepers(),
	Mailer:    Mailer(),
}, 0, time.Second, 10*time.Second)

func AlarmEngine() alarm.Engine {
	return alarmEngine
}

func Protect(fn func()) {
	defer Recover()
	fn()
}

func Recover() {
	alarmEngine.Recover(DeployName())
}

func Alarmf(format string, args ...interface{}) {
	alarmEngine.Alarmf(DeployName()+` `+format, args...)
}

func Alarm(title string) {
	alarmEngine.Alarm(DeployName() + ` ` + title)
}

func AlarmMergeKey(title, content, mergeKey string) {
	alarmEngine.AlarmMergeKey(DeployName()+` `+title, content, mergeKey)
}
