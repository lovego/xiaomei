package alarm

import (
	"fmt"
	"sync"
	"time"

	"github.com/lovego/xiaomei/utils"
)

type Engine interface {
	Recover(prefix string)
	Alarmf(format string, args ...interface{})
	Alarm(title string)
	AlarmMergeKey(title, content, mergeKey string)
}

// engine合并报警邮件，防止在出错高峰，收到大量重复报警邮件，
//	 甚至因为邮件过多导致发送失败、接收失败。
type engine struct {
	sender        Sender
	min, inc, max time.Duration // 发送间隔时间最小值，最大值，增加值
	alarms        map[string]*alarm
	sync.Mutex
}

func NewEngine(sender Sender, min, inc, max time.Duration) Engine {
	return &engine{
		sender: sender,
		min:    min,
		inc:    inc,
		max:    max,
		alarms: make(map[string]*alarm),
	}
}

func (e *engine) Recover(prefix string) {
	if err := recover(); err != nil {
		e.Alarm(fmt.Sprintf("PANIC: %v", err))
	}
}

func (e *engine) Alarmf(format string, args ...interface{}) {
	e.Alarm(fmt.Sprintf(format, args...))
}

// 根据title和调用栈对报警消息进行合并
func (e *engine) Alarm(title string) {
	titleAndStack := title + "\n" + utils.Stack(1)
	content := time.Now().Format(utils.ISO8601) + ` ` + titleAndStack
	println(content)
	e.AlarmMergeKey(title, content, titleAndStack)
}

// 根据mergeKey对报警消息进行合并
func (e *engine) AlarmMergeKey(title, content, mergeKey string) {
	e.Lock()
	a := e.alarms[mergeKey]
	if a == nil {
		a = &alarm{engine: e, interval: e.min, lastSendTime: time.Now()}
		e.alarms[mergeKey] = a
	}
	e.Unlock()
	a.Add(title, content)
}
