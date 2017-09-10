package errs

import (
	"bytes"
	"fmt"
	"runtime"
)

type TraceErr interface {
	Stack() string
	Error() string
}

func Trace(err error) TraceErr {
	if s, ok := err.(trace); ok {
		return s
	} else {
		return trace{err: err.Error(), stack: getStack(1)}
	}
}

func Tracef(format string, args ...interface{}) TraceErr {
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

	callers := make([]uintptr, 32)
	n := runtime.Callers(skip, callers)
	frames := runtime.CallersFrames(callers[:n])
	for {
		if f, ok := frames.Next(); ok {
			fmt.Fprintf(buf, "%s %s:%d (0x%x)\n", f.Function, f.File, f.Line, f.PC)
		} else {
			break
		}
	}
	return buf.String()
}
