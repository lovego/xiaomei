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
var alarm = config.Alarm()

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
		fields[`err`] = fmt.Sprintf("Panic: %v", err)
		fields[`stack`] = utils.Stack(3)
	}

	if line := serializeFields(fields); len(line) > 0 {
		if err != nil {
			theErrorLog.Write(line)
		} else {
			theAccessLog.Write(line)
		}
	}

	if fields[`err`] != nil {
		errStr := fmt.Sprint(fields[`err`])
		errStack := fmt.Sprint(fields[`stack`])
		alarm.Alarm(errStr, formatFields(fields, false), errStr+` `+errStack)
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
