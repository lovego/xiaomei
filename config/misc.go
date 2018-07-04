package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lovego/alarm"
	"github.com/lovego/fs"
	loggerPkg "github.com/lovego/logger"
)

var theAlarm = alarm.New(alarm.MailSender{
	Receivers: Keepers(),
	Mailer:    Mailer(),
}, 0, 5*time.Second, 30*time.Second, alarm.SetPrefix(DeployName()))

var theLogger = loggerPkg.New(os.Stderr)

func init() {
	theLogger.SetAlarm(theAlarm)
	theLogger.SetMachineName()
	theLogger.SetMachineIP()
}

func DevMode() bool {
	return os.Getenv(`GODEV`) == `true`
}

func Alarm() *alarm.Alarm {
	return theAlarm
}

func Logger() *loggerPkg.Logger {
	return theLogger
}

func NewLogger(paths ...string) *loggerPkg.Logger {
	file, err := fs.NewLogFile(filepath.Join(
		append([]string{Root(), `log`}, paths...)...,
	))
	if err != nil {
		Logger().Fatal(err)
	}
	logger := loggerPkg.New(file)
	logger.SetAlarm(Alarm())
	logger.SetMachineName()
	logger.SetMachineIP()
	return logger
}

func Protect(fn func()) {
	defer theLogger.Recover()
	fn()
}
