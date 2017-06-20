package mailer

import (
	"testing"
	"time"

	"github.com/jordan-wright/email"
)

func TestSend(t *testing.T) {
	mailer, err := New(
		`mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh`,
	)
	if err != nil {
		panic(err)
	}

	e := &email.Email{
		To: []string{"侯志良<houzhiliang@retail-tek.com>",
			"侯志良<applejava@qq.com>", "侯志良<bughou@gmail.com>"},
		Subject: "测试",
		HTML:    []byte("<b>超文本!</b>"),
	}
	if err := mailer.Send(e, 10*time.Second); err != nil {
		panic(err)
	}
}
