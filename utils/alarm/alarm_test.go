package alarm

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type testSenderT []string

func (ts *testSenderT) Send(title, content string, count int) {
	*ts = append(*ts, fmt.Sprintf("标题%s[%d] 内容%s", title, count, content))
}

// 外层slice有序，内层slice无序
func (ts testSenderT) Equal(target [][]string) bool {
	now := 0
	for _, slice := range target {
		end := now + len(slice)
		if end > len(ts) || !strSliceEqual(ts[now:end], slice) {
			return false
		}
		now = end
	}
	return now == len(ts)
}

func strSliceEqual(a, b []string) bool {
	m := map[string]bool{}
	for _, str := range a {
		m[str] = true
	}
	for _, str := range b {
		if !m[str] {
			return false
		}
	}
	return true
}

const redundantTime = 10 * time.Millisecond

func TestRecover(t *testing.T) {
	s := &testSenderT{}
	engine := New(`Recover`, s, 0, time.Second, 10*time.Second, nil)
	func() {
		defer engine.Recover()
		panic(`wrong`)
	}()
	time.Sleep(redundantTime) // wait the alarm wait and send goroutine ends.
	fmt.Println(s)
}

func TestPrintf(t *testing.T) {
	s := &testSenderT{}
	engine := New(``, s, 0, time.Second, 10*time.Second, nil)
	engine.Printf("Alarmf: %s", `wrong`)
	time.Sleep(redundantTime) // wait the alarm wait and send goroutine ends.
	fmt.Println(s)
}

func TestAlarm(t *testing.T) {
	s := &testSenderT{}
	engine := New(``, s, 0, time.Second, 10*time.Second, nil)
	engine.Alarm(`Alarm: wrong`)
	time.Sleep(redundantTime) // wait the alarm wait and send goroutine ends.
	fmt.Println(s)
}

func TestDo1(t *testing.T) {
	s := &testSenderT{}
	engine := New(``, s, 0, time.Second, 10*time.Second, nil)
	testDo(engine, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(redundantTime) // wait the test goroutine ends.
	// 首次报警立即发出
	assertEqual(t, s, [][]string{
		{`标题a[1] 内容a`, `标题b[1] 内容b`, `标题c[1] 内容c`},
	})
	// 后续报警等待1秒
	time.Sleep(time.Second + redundantTime) // wait the alarm wait and send goroutine ends.
	assertEqual(t, s, [][]string{
		{`标题a[1] 内容a`, `标题b[1] 内容b`, `标题c[1] 内容c`},
		{`标题a[2] 内容a`, `标题b[3] 内容b`, `标题c[4] 内容c`},
	})
	time.Sleep(2*time.Second + redundantTime)
	testDo(engine, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(redundantTime) // wait the test goroutine ends.
	// 本次报警立即发出
	assertEqual(t, s, [][]string{
		{`标题a[1] 内容a`, `标题b[1] 内容b`, `标题c[1] 内容c`},
		{`标题a[2] 内容a`, `标题b[3] 内容b`, `标题c[4] 内容c`},
		{`标题a[1] 内容a`, `标题b[1] 内容b`, `标题c[1] 内容c`},
	})
}

func TestDo2(t *testing.T) {
	s := &testSenderT{}
	engine := New(``, s, time.Second, time.Second, 10*time.Second, nil)

	testDo(engine, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(redundantTime) // wait the test goroutine ends.
	// 首次报警等待1秒
	time.Sleep(time.Second + redundantTime)
	assertEqual(t, s, [][]string{
		{`标题a[3] 内容a`, `标题b[4] 内容b`, `标题c[5] 内容c`},
	})

	// 先等待1秒
	time.Sleep(time.Second)
	testDo(engine, map[string]int{`a`: 2, `b`: 3, `c`: 1})
	testDo(engine, map[string]int{`a`: 4, `b`: 4, `c`: 7})
	// 本次报警等待2秒，减去之前已经等待的1秒，只需再等待1秒
	time.Sleep(time.Second + redundantTime)
	assertEqual(t, s, [][]string{
		{`标题a[3] 内容a`, `标题b[4] 内容b`, `标题c[5] 内容c`},
		{`标题a[6] 内容a`, `标题b[7] 内容b`, `标题c[8] 内容c`},
	})
}

func testDo(engine *Alarm, groups map[string]int) {
	var wg sync.WaitGroup
	for mergeKey, count := range groups {
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(mergeKey string) {
				defer wg.Done()
				engine.Do(mergeKey, mergeKey, mergeKey)
			}(mergeKey)
		}
	}
	wg.Wait()
}

func assertEqual(t *testing.T, s *testSenderT, expect [][]string) {
	if !s.Equal(expect) {
		t.Errorf("expect: %v\n   got: %v\n", expect, s)
	}
}
