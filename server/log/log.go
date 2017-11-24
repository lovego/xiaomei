package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lovego/errs"
	"github.com/lovego/fs"
	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
)

var isDevMode = config.DevMode()
var theAccessLog, theErrorLog = getLogWriter()
var alarm = config.Alarm()

func getLogWriter() (io.Writer, io.Writer) {
	if isDevMode {
		return os.Stdout, os.Stderr
	}
	path := filepath.Join(config.Root(), `log`, `app`)
	accessLog, err := fs.NewLogFile(path + `.log`)
	if err != nil {
		log.Fatal(err)
	}
	errorLog, err := fs.NewLogFile(path + `.err`)
	if err != nil {
		log.Fatal(err)
	}
	return accessLog, errorLog
}

func Write(req *xiaomei.Request, res *xiaomei.Response, t time.Time, err interface{}) {
	fields := getFields(req, res, t)
	if err != nil {
		fields[`err`] = fmt.Sprintf("Panic: %v", err)
		fields[`stack`] = errs.Stack(3)
	} else if err = fields[`err`]; err != nil {
		if _, ok := err.(string); !ok {
			fields[`err`] = fmt.Sprint(err)
		}
	}

	if line := serializeFields(fields); len(line) > 0 {
		if err != nil {
			theErrorLog.Write(line)
		} else {
			theAccessLog.Write(line)
		}
	}

	if err != nil {
		errStr := fields[`err`].(string)
		errStack, _ := fields[`stack`].(string)
		alarm.Alarm(errStr, formatFields(fields, false), errStr+` `+errStack)
	}
}

func serializeFields(fields map[string]interface{}) []byte {
	if isDevMode {
		return []byte(formatFields(fields, true))
	}
	line, err := json.Marshal(fields)
	if err != nil {
		log.Println(`writeLog:` + err.Error())
	}
	if len(line) > 0 {
		line = append(line, '\n')
	}
	return line
}
