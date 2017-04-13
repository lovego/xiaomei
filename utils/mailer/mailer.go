package mailer

import (
	"log"
	"net/mail"
	"os"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Sender     *mail.Address
	Dialer     *gomail.Dialer
	SendCloser gomail.SendCloser
}

type Message struct {
	*gomail.Message
}

func New(host string, port int, password string, from string) *Mailer {
	if host == `` || password == `` || from == `` {
		return nil
	}
	sender, err := mail.ParseAddress(from)
	if err != nil {
		log.Println(err)
		return nil
	}
	dialer := gomail.NewDialer(host, port, sender.Address, password)
	s, err := dialer.Dial()
	if err != nil {
		panic(err)
	}
	m := Mailer{
		sender, dialer, s,
	}
	return &m
}

func (self *Mailer) Send(msg Message) error {
	return gomail.Send(self.SendCloser, msg.Message)
}

func (self *Mailer) NewMessage(receivers, cc []string, title, body, contentType string) Message {
	m := Message{gomail.NewMessage( /*gomail.SetCharset("UTF-8")*/ )}
	m.setHeaders(self.Sender, parseAdderss(receivers), parseAdderss(cc), title)
	m.setBody(contentType, body)
	return m
}

func (m *Message) setHeaders(from *mail.Address, to, cc []*mail.Address, subject string) {
	m.Message.SetAddressHeader("From", from.Address, from.Name)
	recievers := []string{}
	for _, t := range to {
		recievers = append(recievers, m.Message.FormatAddress(t.Address, t.Name))
	}
	m.SetHeaders(map[string][]string{"To": recievers})
	m.Message.SetHeader("To", recievers...)
	if len(cc) > 0 {
		ccReceivers := []string{}
		for _, c := range cc {
			ccReceivers = append(ccReceivers, m.Message.FormatAddress(c.Address, c.Name))
		}
		m.Message.SetHeader("Cc", ccReceivers...)
	}
	m.Message.SetHeader("Subject", subject)
}

func (m *Message) setBody(contentType, body string) {
	if contentType == `` {
		m.Message.SetBody("text/plain", body)
	} else {
		m.Message.SetBody(contentType, body)
	}
}

func (m *Message) AddAttachs(files ...string) {
	for _, filename := range files {
		if f, err := os.Stat(filename); err != nil {
			log.Printf("WARNING: %s does not exists.\n", filename)
		} else {
			if f.IsDir() {
				log.Printf("WARNING: %s is not a file.\n", filename)
			}
		}
		m.Attach(filename)
	}
}

func parseAdderss(addrs []string) (result []*mail.Address) {
	for _, addr := range addrs {
		if r, err := mail.ParseAddress(addr); err == nil {
			result = append(result, r)
		}
	}
	return
}
