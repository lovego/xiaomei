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

func Stack(err error) StackErr {
	if s, ok := err.(stack); ok {
		return s
	} else {
		return stack{err: err.Error(), stack: getStack(1)}
	}
}

func Stackf(format string, args ...interface{}) StackErr {
	err := fmt.Errorf(format, args)
	return stack{err: err.Error(), stack: getStack(1)}
}

type stack struct {
	err, stack string
}

func (s stack) Stack() string {
	return s.stack
}

func (s stack) Error() string {
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
