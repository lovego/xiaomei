package alarm

import (
	"fmt"
	"time"

	"github.com/jordan-wright/email"
	"github.com/lovego/xiaomei/utils/mailer"
)

type Alarm interface {
	Send(mergedCount int) error
}

type Mail struct {
	Mailer *mailer.Mailer
	Email  *email.Email
}

func (m Mail) Send(count int) error {
	m.Email.Subject += fmt.Sprintf(` [Merged: %d]`, count)
	return m.Mailer.Send(m.Email, time.Minute)
}
