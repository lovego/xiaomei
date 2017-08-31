package alarm

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/lovego/xiaomei/utils/mailer"
)

var testSender = getTestSender()

func getTestSender() Sender {
	theMailer, err := mailer.New(
		`mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh`,
	)
	if err != nil {
		log.Panic(err)
	}

	return MailSender{
		Receivers: []string{"侯志良<houzhiliang@retail-tek.com>"},
		Mailer:    theMailer,
	}
}

func TestAlarm1(t *testing.T) {
	engine := NewEngine(testSender, 0, time.Second, 10*time.Second)
	wg := &sync.WaitGroup{}
	testAlarmMergeKey(engine, wg, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	wg.Wait()
}

func TestAlarm2(t *testing.T) {
	engine := NewEngine(testSender, time.Second, time.Second, 10*time.Second)

	wg := &sync.WaitGroup{}
	testAlarmMergeKey(engine, wg, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(2 * time.Second)
	testAlarmMergeKey(engine, wg, map[string]int{`a`: 2, `b`: 3, `c`: 1})
	testAlarmMergeKey(engine, wg, map[string]int{`a`: 4, `b`: 4, `c`: 7})
	wg.Wait()
}

func testAlarmMergeKey(engine *Engine, wg *sync.WaitGroup, groups map[string]int) {
	for mergeKey, count := range groups {
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(mergeKey string) {
				defer wg.Done()
				engine.AlarmMergeKey(`标题`+mergeKey, `内容`+mergeKey, mergeKey)
			}(mergeKey)
		}
	}
}
