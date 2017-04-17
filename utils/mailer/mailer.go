package mailer

import (
	"net/mail"
	"net/smtp"
	"net/url"
	"strconv"
	"time"

	"github.com/jordan-wright/email"
)

type Mailer struct {
	Sender *mail.Address
	Pool   *email.Pool
}

func New(mailerUrl string) (*Mailer, error) {
	mailer, err := url.Parse(mailerUrl)
	if err != nil {
		return nil, err
	}
	query := mailer.Query()

	poolSize := 10
	if size := query.Get(`poolSize`); size != `` {
		if sizeInt, err := strconv.Atoi(size); err != nil {
			return nil, err
		} else if sizeInt > 0 {
			poolSize = sizeInt
		}
	}

	sender, err := mail.ParseAddress(query.Get(`user`))
	if err != nil {
		return nil, err
	}

	pool := email.NewPool(
		mailer.Host, poolSize,
		smtp.PlainAuth(``, sender.Address, query.Get(`pass`), mailer.Hostname()),
	)

	return &Mailer{sender, pool}, nil
}

func (m *Mailer) Send(e *email.Email, timeout time.Duration) error {
	if m == nil {
		return nil
	}
	if e.From == `` && m.Sender != nil {
		e.From = m.Sender.String()
	}
	return m.Pool.Send(e, timeout)
}
