package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

var isDevMode = config.DevMode()
var theAccessLog, theErrorLog = getLogWriter()

func getLogWriter() (io.Writer, io.Writer) {
	if isDevMode {
		return os.Stdout, os.Stderr
	}
	path := filepath.Join(config.Root(), `log`, `app`)
	return fs.NewLogFile(path + `.log`), fs.NewLogFile(path + `.err`)
}

func Write(req *xiaomei.Request, res *xiaomei.Response, t time.Time, err interface{}) {
	fields := getFields(req, res, t)

	if err != nil {
		errStr := fmt.Sprint(err)
		errStack := string(utils.Stack(3))
		fields[`err`] = errStr
		fields[`stack`] = errStack

		if !isDevMode {
			go config.Alarm(`500错误`, formatFields(fields, false), errStr+` `+errStack)
		}
	} else if err = fields[`err`]; err != nil && !isDevMode {
		errStr, ok := err.(string)
		if !ok {
			errStr = fmt.Sprint(err)
		}
		go config.Alarm(`错误`, formatFields(fields, false), errStr)
	}

	if line := serializeFields(fields); len(line) > 0 {
		if err != nil {
			theErrorLog.Write(line)
		} else {
			theAccessLog.Write(line)
		}
	}
}

func serializeFields(fields map[string]interface{}) []byte {
	if isDevMode {
		return []byte(formatFields(fields, true))
	}
	line, err := json.Marshal(fields)
	if err != nil {
		utils.Log(`writeLog:` + err.Error())
	}
	if len(line) > 0 {
		line = append(line, '\n')
	}
	return line
}
