package mailer

import (
	"sync"
	"testing"
	"time"

	"github.com/jordan-wright/email"
)

func TestSend(t *testing.T) {
	mailer, err := New(
		`mailer://imap.exmail.qq.com:25/?user=大数据<name@qq.com>&pass=password`,
	)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for _, to := range []string{
		"侯志良<houzhiliang@retail-tek.com>",
		"侯志良<applejava@qq.com>",
		"侯志良<bughou@gmail.com>",
	} {
		wg.Add(1)
		go func(to string) {
			defer wg.Done()

			e := &email.Email{
				To:      []string{to},
				Subject: "测试-" + to,
				HTML:    []byte("<b>超文本!</b>"),
			}
			if err := mailer.Send(e, time.Minute); err != nil {
				panic(err)
			}
		}(to)
	}
	wg.Wait()
}
