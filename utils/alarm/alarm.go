package alarm

import (
	"sync"
	"time"
)

// Alarm合并报警邮件，防止在出错高峰，收到大量重复报警邮件，
//	 甚至因为邮件过多导致发送失败、接收失败。
type Alarm struct {
	prefix        string
	sender        Sender
	min, inc, max time.Duration // 发送间隔时间最小值，最大值，增加值
	sync.Mutex
	alarms map[string]*alarm
}

func New(
	prefix string, sender Sender, min, inc, max time.Duration,
) *Alarm {
	return &Alarm{
		prefix: prefix,
		sender: sender,
		min:    min,
		inc:    inc,
		max:    max,
		alarms: make(map[string]*alarm),
	}
}

func (alm *Alarm) Send(title, content string) {
	if alm == nil {
		return
	}
	alm.sender.Send(alm.prefix+` `+title, content, 1)
}

func (alm *Alarm) Alarm(title, content, mergeKey string) {
	if alm == nil {
		return
	}
	alm.Lock()
	// 根据mergeKey对报警消息进行合并
	a := alm.alarms[mergeKey]
	if a == nil {
		a = &alarm{Alarm: alm, interval: alm.min, lastSendTime: time.Now()}
		alm.alarms[mergeKey] = a
	}
	alm.Unlock()
	if alm.prefix != `` {
		title = alm.prefix + ` ` + title
	}
	a.Add(title, content)
}
