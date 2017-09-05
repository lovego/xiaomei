package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lovego/xiaomei/utils"
)

type Logger struct {
	prefix string
	writer io.Writer
	alarm  Alarm
}

type Alarm interface {
	Send(title, content string)
	Alarm(title, content, mergeKey string)
}

const timeFormat = `2006/01/02 15:04:05`

func New(prefix string, writer io.Writer, alarm Alarm) *Logger {
	if writer == nil {
		writer = os.Stderr
	}
	return &Logger{
		prefix: prefix, writer: writer, alarm: alarm,
	}
}

func (log *Logger) Print(args ...interface{}) {
	log.writer.Write([]byte(log.output(fmt.Sprint(args...))))
}

func (log *Logger) Printf(format string, args ...interface{}) {
	log.writer.Write([]byte(log.output(fmt.Sprintf(format, args...))))
}

func (log *Logger) Println(args ...interface{}) {
	log.writer.Write([]byte(log.output(fmt.Sprintln(args...))))
}

func (log *Logger) output(s string) string {
	if len(s) == 0 || s[len(s)-1] != '\n' {
		s += "\n"
	}
	return log.prefix + time.Now().Format(timeFormat) + ` ` + s
}

func (log *Logger) Error(args ...interface{}) {
	log.doAlarm(fmt.Sprint(args...))
}

func (log *Logger) Errorf(format string, args ...interface{}) {
	log.doAlarm(fmt.Sprintf(format, args...))
}

func (log *Logger) Errorln(args ...interface{}) {
	log.doAlarm(fmt.Sprintln(args...))
}

func (log *Logger) doAlarm(title string) {
	stack := utils.Stack(3)
	content := log.output(title) + stack
	log.writer.Write([]byte(content))
	title = log.prefix + title
	mergeKey := title + "\n" + stack // 根据title和调用栈对报警消息进行合并
	log.alarm.Alarm(title, content, mergeKey)
}

func (log *Logger) Recover() {
	if err := recover(); err != nil {
		log.doAlarm(fmt.Sprintf("PANIC: %v\n", err))
	}
}

func (log *Logger) Fatal(args ...interface{}) {
	log.doExit(fmt.Sprint(args...))
}

func (log *Logger) Fatalf(format string, args ...interface{}) {
	log.doExit(fmt.Sprintf(format, args...))
}

func (log *Logger) Fatalln(args ...interface{}) {
	log.doExit(fmt.Sprintln(args...))
}

func (log *Logger) doExit(title string) {
	stack := utils.Stack(3)
	content := log.output(title) + stack
	log.writer.Write([]byte(content))
	title = log.prefix + title
	log.alarm.Send(title, content)
	os.Exit(1)
}

func (log *Logger) Panic(args ...interface{}) {
	log.doPanic(fmt.Sprint(args...))
}

func (log *Logger) Panicf(format string, args ...interface{}) {
	log.doPanic(fmt.Sprintf(format, args...))
}

func (log *Logger) Panicln(args ...interface{}) {
	log.doPanic(fmt.Sprintln(args...))
}

func (log *Logger) doPanic(title string) {
	stack := utils.Stack(3)
	titleLine := log.output(title)

	content := titleLine + stack
	if log.writer != os.Stderr {
		log.writer.Write([]byte(content))
	}
	title = log.prefix + title
	log.alarm.Send(title, content)

	os.Stderr.Write([]byte(titleLine))
	panic(title)
}
