package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/server/xm"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

var accessLog, errLog = setupLogger()

func setupLogger() (*os.File, *os.File) {
	if config.DevMode() {
		return os.Stdout, os.Stderr
	} else {
		return fs.OpenAppend(filepath.Join(config.Root(), `log/app.log`)),
			fs.OpenAppend(filepath.Join(config.Root(), `log/app.err`))
	}
}

func writeLog(
	req *xm.Request, res *xm.Response, t time.Time, hasErr bool, errStr, stack string,
) []byte {
	line, err := json.Marshal(getLogFields(req, res, t, hasErr, errStr, stack))
	if err != nil {
		utils.Log(`writeLog:` + err.Error())
		return nil
	}
	line = append(line, '\n')
	if hasErr {
		errLog.Write(line)
	} else {
		accessLog.Write(line)
	}
	return line
}

func getLogFields(
	req *xm.Request, res *xm.Response, t time.Time, hasErr bool, errStr, stack string,
) map[string]interface{} {
	var sess interface{}
	req.Session(&sess)
	m := map[string]interface{}{
		`at`: t.Format(utils.ISO8601), `time`: fmt.Sprintf(`%.6f`, time.Since(t).Seconds()),
		`host`: req.Host, `method`: req.Method, `path`: req.URL.Path, `query`: req.URL.RawQuery,
		`status`: res.Status(), `req_body`: req.ContentLength, `res_body`: res.Size(),
		`session`: sess, `ip`: req.ClientAddr(),
		`refer`: req.Referer(), `agent`: req.UserAgent(), `proto`: req.Proto,
	}
	if hasErr {
		m[`err`] = errStr
		m[`stack`] = stack
	}
	for k, v := range req.Log {
		m[k] = v
	}
	return m
}
