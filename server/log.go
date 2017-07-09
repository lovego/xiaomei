package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

var accessLog, errorLog = setupLogger()
var isDevMode = config.DevMode()

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
	errorLogPath := filepath.Join(config.Root(), `log/app.err`)
	if errorLog, err = fs.OpenAppend(errorLogPath); err != nil {
		utils.Logf(`open appserver error log %s failed: %v`, errorLogPath, err)
		os.Exit(1)
	}
	return
}

func writeLog(
	req *xiaomei.Request, res *xiaomei.Response, t time.Time, hasErr bool, errStr, stack string,
) []byte {
	line, err := json.Marshal(getLogFields(req, res, t, hasErr, errStr, stack))
	if err != nil {
		utils.Log(`writeLog:` + err.Error())
		return nil
	}
	line = append(line, '\n')
	if hasErr {
		errorLog.Write(line)
		if isDevMode {
			errorLog.WriteString(errStr + "\n" + stack)
		}
	} else {
		accessLog.Write(line)
	}
	return line
}

func getLogFields(
	req *xiaomei.Request, res *xiaomei.Response, t time.Time, hasErr bool, errStr, stack string,
) map[string]interface{} {
	var sess interface{}
	req.Session(&sess)
	m := map[string]interface{}{
		`at`: t.Format(utils.ISO8601), `duration`: fmt.Sprintf(`%.6f`, time.Since(t).Seconds()),
		`host`: req.Host, `method`: req.Method, `path`: req.URL.Path, `query`: req.URL.RawQuery,
		`status`: res.Status(), `req_body`: req.ContentLength, `res_body`: res.Size(),
		`ip`:    req.ClientAddr(),
		`refer`: req.Referer(), `agent`: req.UserAgent(), `proto`: req.Proto,
	}
	if sess != nil {
		m[`session`] = sess
	}
	if hasErr && !isDevMode {
		m[`err`] = errStr
		m[`stack`] = stack
	}
	for k, v := range req.Log {
		m[k] = v
	}
	return m
}
