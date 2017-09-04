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

// Alarm合并报警邮件，防止在出错高峰，收到大量重复报警邮件，
//	 甚至因为邮件过多导致发送失败、接收失败。
type Alarm struct {
	prefix        string
	sender        Sender
	min, inc, max time.Duration // 发送间隔时间最小值，最大值，增加值
	sync.Mutex
	alarms map[string]*alarm
	writer io.Writer
}

func New(
	prefix string, sender Sender, min, inc, max time.Duration, writer io.Writer,
) *Alarm {
	if writer == nil {
		writer = os.Stderr
	}
	return &Alarm{
		prefix: prefix,
		sender: sender,
		min:    min,
		inc:    inc,
		max:    max,
		alarms: make(map[string]*alarm),
		writer: writer,
	}
}

func (alm *Alarm) DupWriter(writer io.Writer) *Alarm {
	dup := *alm
	dup.writer = writer
	return &dup
}

func (alm *Alarm) Log(args ...interface{}) {
	alm.writer.Write([]byte(time.Now().Format(timeFormat) + ` ` + fmt.Sprint(args...)))
}

func (alm *Alarm) Logf(format string, args ...interface{}) {
	alm.writer.Write([]byte(time.Now().Format(timeFormat) + ` ` + fmt.Sprintf(format, args...)))
}

func (alm *Alarm) Recover() {
	if err := recover(); err != nil {
		alm.Printf("PANIC: %v", err)
	}
}

func (alm *Alarm) Fatal(args ...interface{}) {
	title := fmt.Sprint(args...)
	content, _ := alm.getContentMergeKey(title)
	alm.writer.Write([]byte(content))
	alm.sender.Send(alm.prefix+` `+title, content, 0)
	os.Exit(1)
}

func (alm *Alarm) Fatalf(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	content, _ := alm.getContentMergeKey(title)
	alm.writer.Write([]byte(content))
	alm.sender.Send(alm.prefix+` `+title, content, 1)
	os.Exit(1)
}

func (alm *Alarm) Panic(args ...interface{}) {
	title := fmt.Sprint(args...)
	content, _ := alm.getContentMergeKey(title)
	alm.writer.Write([]byte(content))
	alm.sender.Send(alm.prefix+` `+title, content, 1)
	panic(content)
}

func (alm *Alarm) Panicf(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	content, _ := alm.getContentMergeKey(title)
	alm.writer.Write([]byte(content))
	alm.sender.Send(alm.prefix+` `+title, content, 1)
	panic(content)
}

func (alm *Alarm) Print(args ...interface{}) {
	title := fmt.Sprint(args...)
	content, mergeKey := alm.getContentMergeKey(title)
	alm.writer.Write([]byte(content))
	alm.Do(title, content, mergeKey)
}

func (alm *Alarm) Printf(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	content, mergeKey := alm.getContentMergeKey(title)
	alm.writer.Write([]byte(content))
	alm.Do(title, content, mergeKey)
}

func (alm *Alarm) Alarm(title string) {
	content, mergeKey := alm.getContentMergeKey(title)
	alm.Do(title, content, mergeKey)
}

func (alm *Alarm) Do(title, content, mergeKey string) {
	alm.Lock()
	// 根据mergeKey对报警消息进行合并
	a := alm.alarms[mergeKey]
	if a == nil {
		a = &alarm{Alarm: alm, interval: alm.min, lastSendTime: time.Now()}
		alm.alarms[mergeKey] = a
	}
	alm.Unlock()
	a.Add(alm.prefix+` `+title, content)
}

func (alm *Alarm) getContentMergeKey(title string) (string, string) {
	// 根据title和调用栈对报警消息进行合并
	mergeKey := title + "\n" + utils.Stack(3)
	content := time.Now().Format(timeFormat) + ` ` + mergeKey
	return content, mergeKey
}
