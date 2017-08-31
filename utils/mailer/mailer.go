package mailer

import (
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/lovego/xiaomei/utils"
	"io"
)

type Mailer struct {
	Sender   *mail.Address
	Pool     *email.Pool
	poolSize int
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

	return &Mailer{sender, pool, poolSize}, nil
}

func (m *Mailer) Send(e *email.Email, timeout time.Duration) (err error) {
	if m == nil || e == nil || len(e.To) == 0 {
		return nil
	}
	if e.From == `` && m.Sender != nil {
		e.From = m.Sender.String()
	}
	setupAddrsHeaders(e)

	// 如果是io.EOF错误,可能是由于pool连接被关闭造成, 重试
	for i := 0; i < m.poolSize; i++ {
		err = m.Pool.Send(e, timeout)
		if err == nil || err != io.EOF {
			return
		}
	}
	return
}

func setupAddrsHeaders(e *email.Email) {
	if e.Headers == nil {
		e.Headers = make(textproto.MIMEHeader)
	}
	if len(e.From) > 0 {
		e.Headers.Set(`From`, quoteAddr(e.From))
	}
	if len(e.To) > 0 {
		e.Headers.Set(`To`, makeAddrsHeader(e.To))
	}
	if len(e.Cc) > 0 {
		e.Headers.Set(`Cc`, makeAddrsHeader(e.Cc))
	}
}

func makeAddrsHeader(addrs []string) string {
	var result []string
	for _, addr := range addrs {
		result = append(result, quoteAddr(addr))
	}
	return strings.Join(result, `, `)
}

func quoteAddr(addr string) string {
	if address, err := mail.ParseAddress(addr); err != nil {
		utils.Log(err)
		return ``
	} else {
		return fmt.Sprintf(`%s <%s>`, mime.QEncoding.Encode(`UTF-8`, address.Name), address.Address)
	}
}
