package alarm

import (
	"fmt"
	"log"
	"time"

	"github.com/jordan-wright/email"
	"github.com/lovego/xiaomei/utils/mailer"
)

type Sender interface {
	Send(title, content string, count int)
}

type MailSender struct {
	Receivers []string
	Mailer    *mailer.Mailer
}

func (m MailSender) Send(title, content string, count int) {
	if len(m.Receivers) == 0 {
		return
	}
	if count > 1 {
		title += fmt.Sprintf(` [Merged: %d]`, count)
	}

	err := m.Mailer.Send(&email.Email{
		To:      m.Receivers,
		Subject: title,
		Text:    []byte(content),
	}, time.Minute)

	if err != nil {
		log.Printf("send alarm mail failed: %v", err)
	}
}
