package alarm

import (
	"sync"
	"time"
)

type mergedAlarm struct {
	sync.Mutex
	alarm        Alarm
	count        int
	lastSendTime time.Time
	interval     time.Duration
}

func (ma *mergedAlarm) Add(a Alarm, min, inc, max time.Duration) error {
	ma.Lock()
	ma.count++
	count := ma.count
	ma.Unlock()

	if count == 1 { // 发送间隔内的首次报警
		ma.alarm = a
		ma.Wait(min)
		return ma.Send(inc, max)
	} else {
		return nil
	}
}

func (ma *mergedAlarm) Wait(min time.Duration) {
	if ma.interval <= 0 {
		return
	}
	if gap := ma.interval - time.Since(ma.lastSendTime); gap > 0 {
		time.Sleep(gap)
	} else {
		// 本次报警超过了间隔时间, 重置间隔时间为min
		ma.interval = min
		if min > 0 {
			time.Sleep(min)
		}
	}
}

func (ma *mergedAlarm) Send(inc, max time.Duration) error {
	ma.Lock()
	alarm := ma.alarm
	count := ma.count
	ma.alarm = nil
	ma.count = 0
	ma.lastSendTime = time.Now()
	ma.interval += inc // 每发送一次，间隔时间增加inc，直到max。
	if ma.interval > max {
		ma.interval = max
	}
	ma.Unlock()

	return alarm.Send(count)
}
