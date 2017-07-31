package alarm

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/jordan-wright/email"
	"github.com/lovego/xiaomei/utils/mailer"
)

var testMailer = getTestMailer()

func getTestMailer() *mailer.Mailer {
	theMailer, err := mailer.New(
		`mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh`,
	)
	if err != nil {
		log.Panic(err)
	}
	return theMailer
}

func TestAlarm1(t *testing.T) {
	engine := NewEngine(0, time.Second, 10*time.Second)
	wg := &sync.WaitGroup{}
	testAlarm(engine, wg, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	wg.Wait()
}

func TestAlarm2(t *testing.T) {
	engine := NewEngine(time.Second, time.Second, 10*time.Second)

	wg := &sync.WaitGroup{}
	testAlarm(engine, wg, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(time.Second)
	testAlarm(engine, wg, map[string]int{`a`: 2, `b`: 3, `c`: 1})
	testAlarm(engine, wg, map[string]int{`a`: 4, `b`: 4, `c`: 7})
	wg.Wait()
}

func testAlarm(engine *Engine, wg *sync.WaitGroup, groups map[string]int) {
	for mergeKey, count := range groups {
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(mergeKey string) {
				defer wg.Done()

				if err := engine.Alarm(makeTestAlarm(mergeKey), mergeKey); err != nil {
					log.Panic(err)
				}
			}(mergeKey)
		}
	}
}

func makeTestAlarm(mergeKey string) Alarm {
	return Mail{
		Mailer: testMailer,
		Email: &email.Email{
			To:      []string{"侯志良<houzhiliang@retail-tek.com>"},
			Subject: "测试-" + mergeKey,
			HTML:    []byte("<b>超文本!</b>"),
		},
	}
}
