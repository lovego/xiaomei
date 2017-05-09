package mailer

import (
	"sync"
	"testing"
	"time"

	"github.com/jordan-wright/email"
)

func TestAlarm(t *testing.T) {
	mailer, err := New(
		`mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=gtrqcegfgtqwebga`,
	)
	if err != nil {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	testAlarm(mailer, wg, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(4 * time.Second)
	testAlarm(mailer, wg, map[string]int{`a`: 2, `b`: 3, `c`: 1})
	testAlarm(mailer, wg, map[string]int{`a`: 4, `b`: 4, `c`: 7})
	wg.Wait()
}

func testAlarm(mailer *Mailer, wg *sync.WaitGroup, groups map[string]int) {
	for mergeKey, count := range groups {
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(mergeKey string) {
				defer wg.Done()

				e := &email.Email{
					To:      []string{"侯志良<houzhiliang@retail-tek.com>"},
					Subject: "测试-" + mergeKey,
					HTML:    []byte("<b>超文本!</b>"),
				}
				if err := mailer.Alarm(e, mergeKey); err != nil {
					panic(err)
				}
			}(mergeKey)
		}
	}
}
