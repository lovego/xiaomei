package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

const ISO8601 = `2006-01-02T15:04:05Z0700`

func Log(args ...interface{}) {
	args = append([]interface{}{time.Now().Format(ISO8601)}, args...)
	fmt.Println(args...)
}

func Logf(format string, args ...interface{}) {
	format = `%s ` + format + "\n"
	args = append([]interface{}{time.Now().Format(ISO8601)}, args...)
	fmt.Printf(format, args...)
}

func Flog(w io.Writer, args ...interface{}) {
	args = append([]interface{}{time.Now().Format(ISO8601)}, args...)
	fmt.Fprintln(w, args...)
}

func FLogf(w io.Writer, format string, args ...interface{}) {
	format = `%s ` + format + "\n"
	args = append([]interface{}{time.Now().Format(ISO8601)}, args...)
	fmt.Fprintf(w, format, args...)
}

func Stack(skip int) string {
	buf := new(bytes.Buffer)

	callers := make([]uintptr, 32)
	n := runtime.Callers(skip, callers)
	frames := runtime.CallersFrames(callers[:n])
	for {
		if f, ok := frames.Next(); ok {
			fmt.Fprintf(buf, "%s\n\t%s:%d (0x%x)\n", f.Function, f.File, f.Line, f.PC)
		} else {
			break
		}
	}
	return buf.String()
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
