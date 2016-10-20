package xm

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
)

func Stack(skip int) []byte {
	buf := new(bytes.Buffer)
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.Bytes()
}

func Protect(fn func()) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("PANIC: %s\n%s", err, Stack(4))
		}
	}()
	fn()
}
