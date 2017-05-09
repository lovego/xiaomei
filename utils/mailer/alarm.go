package mailer

import (
	"fmt"
	"sync"
	"time"

	"github.com/jordan-wright/email"
)

type alarmEmail struct {
	email *email.Email
	count int
}

var alarmEmails = struct {
	sync.Mutex
	m map[string]*alarmEmail
}{
	m: make(map[string]*alarmEmail),
}

func (m *Mailer) Alarm(e *email.Email, mergeKey string) error {
	if mergeKey == `` {
		return m.Send(e, time.Minute)
	}

	alarmEmails.Lock()
	ae := alarmEmails.m[mergeKey]
	if ae == nil {
		ae = &alarmEmail{email: e, count: 1}
		alarmEmails.m[mergeKey] = ae
	} else {
		ae.count++
	}
	count := ae.count
	alarmEmails.Unlock()

	if count != 1 {
		return nil
	}
	time.Sleep(3 * time.Second)

	alarmEmails.Lock()
	delete(alarmEmails.m, mergeKey)
	alarmEmails.Unlock()

	ae.email.Subject += fmt.Sprintf(` [Merged: %d]`, ae.count)
	return m.Send(ae.email, time.Minute)
}
