package mailer

import (
	"testing"
)

func TestSend(t *testing.T) {
	m := New(`imap.exmail.qq.com`, 25, `RetailBigdata1234`, `data-system<data-system@retail-tek.com>`)
	msg := m.NewMessage([]string{`chen<984258288@qq.com>`}, nil, `title`, `body`, ``)
	err := m.Send(msg)
	if err != nil {
		panic(err)
	}
}
