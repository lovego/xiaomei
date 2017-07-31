package alarm

import (
	"sync"
	"time"
)

// Engine合并报警邮件，防止在出错高峰，收到大量重复报警邮件，
//	 甚至因为邮件过多导致发送失败、接收失败。
type Engine struct {
	min, inc, max time.Duration // 发送间隔时间最小值，最大值，增加值
	sync.Mutex
	alarms map[string]*mergedAlarm
}

func NewEngine(min, inc, max time.Duration) *Engine {
	return &Engine{
		min:    min,
		inc:    inc,
		max:    max,
		alarms: make(map[string]*mergedAlarm),
	}
}

/*
	根据mergeKey对报警进行合并
*/
func (e *Engine) Alarm(a Alarm, mergeKey string) error {
	if mergeKey == `` {
		return a.Send(0)
	}

	e.Lock()
	ma := e.alarms[mergeKey]
	if ma == nil {
		ma = &mergedAlarm{interval: e.min, lastSendTime: time.Now()}
		e.alarms[mergeKey] = ma
	}
	e.Unlock()
	return ma.Add(a, e.min, e.inc, e.max)
}
