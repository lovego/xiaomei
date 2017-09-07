package errs

import (
	"bytes"
	"fmt"
	"runtime"
)

type StackErr interface {
	Stack() string
	Error() string
}

func Trace(err error) StackErr {
	if s, ok := err.(trace); ok {
		return s
	} else {
		return trace{err: err.Error(), stack: getStack(1)}
	}
}

func Stackf(format string, args ...interface{}) StackErr {
	return trace{err: fmt.Sprintf(format, args...), stack: getStack(1)}
}

type trace struct {
	err, stack string
}

func (s trace) Stack() string {
	return s.stack
}

func (s trace) Error() string {
	return s.err
}

func getStack(skip int) string {
	buf := new(bytes.Buffer)
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}
