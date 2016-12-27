package mailer

import (
	"testing"
)

var mailer = New(
	`imap.mail.qq.com`, `25`, `小美<xiaomei-go@qq.com>`, `abc123`,
)

func TestSend(t *testing.T) {
	mailer.Send(&Message{
		Receivers: []string{`小美<xiaomei-go@qq.com>`},
		Title:     `test 标题`,
		Body:      `test 内容`,
	})
}

func TestSendWithAttaches(t *testing.T) {
	files := map[string]string{
		`test1.txt`: `test1`,
		`test2.txt`: `test2`,
		`测试3.txt`:   `测试3`,
	}
	mailer.Send(&Message{
		Receivers: []string{`小美<xiaomei-go@qq.com>`},
		Title:     `test 标题`,
		Body:      `test 内容`,
		Attaches:  files,
	})
}
