package mailer

import (
	"testing"
)

var mailer = New(
	`imap.exmail.qq.com`, `25`, `data-system<data-system@retail-tek.com>`, `RetailBigdata1234`,
)

func TestSend(t *testing.T) {
	mailer.Send(
		[]string{`侯志良<houzhiliang@retail-tek.com>`}, `test 标题`, `test 内容`,
	)
}

func TestSendWithAttaches(t *testing.T) {
	files := map[string]string{
		`test1.txt`: `test1`,
		`test2.txt`: `test2`,
		`测试3.txt`:   `test3`,
	}
	mailer.SendWithAttaches(
		[]string{`侯志良<houzhiliang@retail-tek.com>`}, `test 标题`, `test 内容`, files,
	)
}
