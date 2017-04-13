package mailer

import (
	"testing"
)

func TestSend(t *testing.T) {
	m, err := New(`imap.exmail.qq.com`, 25, `RetailBigdata1234`, `data-system<data-system@retail-tek.com>`)
	if err != nil {
		panic(err)
	}
	msg := m.NewMessage([]string{`chen<984258288@qq.com>`}, nil, ``, ``, ``)
	err = m.Send(msg)
	if err != nil {
		panic(err)
	}
}
