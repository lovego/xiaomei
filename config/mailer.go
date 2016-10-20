package config

import (
	"github.com/bughou-go/xiaomei/utils/mailer"
)

var Mailer *mailer.Mailer

func setupMailer() {
	m := Data.Mailer
	if m.Host == `` || m.Port == `` || m.Sender == `` {
		return
	}
	Mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
}

func AlarmMail(title, body string) {
	title = Data.DeployName + ` ` + title
	Mailer.Send(Data.AlarmReceivers, title, body)
}
