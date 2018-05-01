package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

func Write(req *xiaomei.Request, res *xiaomei.Response, panicError interface{}) {
	fields := getFields(req, res)
	if panicError != nil {
		fields.Error = fmt.Sprintf("Panic: %v", panicError)
		fields.Stack = errs.Stack(3)
	}
	if line := serializeFields(fields); len(line) > 0 {
		if fields.Error != "" {
			theErrorLog.Write(line)
		} else {
			theAccessLog.Write(line)
		}
	}

	if fields.Error != "" {
		alarm.Alarm(fields.Error, formatFields(fields, false), fields.Error+` `+fields.Stack)
	}
}

func serializeFields(fields *logFields) []byte {
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
