package alarm

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/lovego/xiaomei/utils"
)

const timeFormat = `2006-01-02 15:04:05`

// Engine合并报警邮件，防止在出错高峰，收到大量重复报警邮件，
//	 甚至因为邮件过多导致发送失败、接收失败。
type Engine struct {
	prefix        string
	sender        Sender
	min, inc, max time.Duration // 发送间隔时间最小值，最大值，增加值
	sync.Mutex
	alarms map[string]*alarm
	writer io.Writer
}

func NewEngine(
	prefix string, sender Sender, min, inc, max time.Duration, writer io.Writer,
) *Engine {
	if writer == nil {
		writer = os.Stderr
	}
	return &Engine{
		prefix: prefix,
		sender: sender,
		min:    min,
		inc:    inc,
		max:    max,
		alarms: make(map[string]*alarm),
		writer: writer,
	}
}

func (e *Engine) Log(args ...interface{}) {
	e.writer.Write([]byte(time.Now().Format(timeFormat) + ` ` + fmt.Sprint(args...)))
}

func (e *Engine) Logf(format string, args ...interface{}) {
	e.writer.Write([]byte(time.Now().Format(timeFormat) + ` ` + fmt.Sprintf(format, args...)))
}

func (e *Engine) Recover() {
	if err := recover(); err != nil {
		e.Printf("PANIC: %v", err)
	}
}

func (e *Engine) Fatal(args ...interface{}) {
	title := fmt.Sprint(args...)
	content, _ := e.getContentMergeKey(title)
	e.writer.Write([]byte(content))
	e.sender.Send(e.prefix+` `+title, content, 0)
	os.Exit(1)
}

func (e *Engine) Fatalf(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	content, _ := e.getContentMergeKey(title)
	e.writer.Write([]byte(content))
	e.sender.Send(e.prefix+` `+title, content, 1)
	os.Exit(1)
}

func (e *Engine) Panic(args ...interface{}) {
	title := fmt.Sprint(args...)
	content, _ := e.getContentMergeKey(title)
	e.writer.Write([]byte(content))
	e.sender.Send(e.prefix+` `+title, content, 1)
	panic(content)
}

func (e *Engine) Panicf(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	content, _ := e.getContentMergeKey(title)
	e.writer.Write([]byte(content))
	e.sender.Send(e.prefix+` `+title, content, 1)
	panic(content)
}

func (e *Engine) Print(args ...interface{}) {
	title := fmt.Sprint(args...)
	content, mergeKey := e.getContentMergeKey(title)
	e.writer.Write([]byte(content))
	e.Do(title, content, mergeKey)
}

func (e *Engine) Printf(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	content, mergeKey := e.getContentMergeKey(title)
	e.writer.Write([]byte(content))
	e.Do(title, content, mergeKey)
}

func (e *Engine) Alarm(title string) {
	content, mergeKey := e.getContentMergeKey(title)
	e.Do(title, content, mergeKey)
}

func (e *Engine) Do(title, content, mergeKey string) {
	e.Lock()
	// 根据mergeKey对报警消息进行合并
	a := e.alarms[mergeKey]
	if a == nil {
		a = &alarm{engine: e, interval: e.min, lastSendTime: time.Now()}
		e.alarms[mergeKey] = a
	}
	e.Unlock()
	a.Add(e.prefix+` `+title, content)
}

func (e *Engine) getContentMergeKey(title string) (string, string) {
	// 根据title和调用栈对报警消息进行合并
	mergeKey := title + "\n" + utils.Stack(3)
	content := time.Now().Format(timeFormat) + ` ` + mergeKey
	return content, mergeKey
}
