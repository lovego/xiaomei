package utils

import (
	"bytes"
	"encoding/json"
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

func PrintJson(v interface{}) {
	data, err := json.MarshalIndent(v, ``, `  `)
	fmt.Println(string(data), err)
}
