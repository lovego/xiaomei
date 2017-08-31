package alarm

import (
	"sync"
	"time"
)

type alarm struct {
	sync.Mutex
	title, content string
	count          int
	lastSendTime   time.Time
	interval       time.Duration
	engine         *engine
}

func (a *alarm) Add(title, content string) {
	a.Lock()
	a.count++
	count := a.count
	a.Unlock()

	if count == 1 { // 发送间隔内的首次报警
		a.title, a.content = title, content
		go func() {
			a.Wait()
			a.Send()
		}()
	}
}

func (a *alarm) Wait() {
	if a.interval <= 0 {
		return
	}
	if gap := a.interval - time.Since(a.lastSendTime); gap > 0 {
		time.Sleep(gap)
	} else {
		// 本次报警超过了间隔时间, 重置间隔时间为min
		a.interval = a.engine.min
		if a.engine.min > 0 {
			time.Sleep(a.engine.min)
		}
	}
}

func (a *alarm) Send() {
	a.Lock()
	title, content, count := a.title, a.content, a.count
	a.title, a.content, a.count = ``, ``, 0
	a.lastSendTime = time.Now()
	a.interval += a.engine.inc // 每发送一次，间隔时间增加inc，直到max。
	if a.interval > a.engine.max {
		a.interval = a.engine.max
	}
	a.Unlock()

	a.engine.sender.Send(title, content, count)
}
