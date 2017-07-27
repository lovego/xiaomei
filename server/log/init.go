package log

import (
	"os"
	"path/filepath"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

var isDevMode = config.DevMode()
var accessLog, errorLog = setupLogger()

func setupLogger() (accessLog, errorLog *os.File) {
	if isDevMode {
		return os.Stdout, os.Stderr
	}
	var err error
	logDir := filepath.Join(config.Root(), `log`)
	if err = os.MkdirAll(logDir, 0775); err != nil {
		utils.Logf(`open appserver log dir %s failed: %v`, logDir, err)
		os.Exit(1)
	}
	accessLogPath := filepath.Join(logDir, `app.log`)
	if accessLog, err = fs.OpenAppend(accessLogPath); err != nil {
		utils.Logf(`open appserver access log %s failed: %v`, accessLogPath, err)
		os.Exit(1)
	}
	errorLogPath := filepath.Join(logDir, `app.err`)
	if errorLog, err = fs.OpenAppend(errorLogPath); err != nil {
		utils.Logf(`open appserver error log %s failed: %v`, errorLogPath, err)
		os.Exit(1)
	}
	return
}
