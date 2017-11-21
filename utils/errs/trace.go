package errs

import (
	"bytes"
	"fmt"
	"runtime"
)

type TraceErr struct {
	err, stack string
}

func Trace(err error) TraceErr {
	if trace, ok := err.(TraceErr); ok {
		return trace
	} else {
		return TraceErr{err: err.Error(), stack: getStack(1)}
	}
}

func Tracef(format string, args ...interface{}) TraceErr {
	return TraceErr{err: fmt.Sprintf(format, args...), stack: getStack(1)}
}

func (s TraceErr) Stack() string {
	return s.stack
}

func (s TraceErr) Error() string {
	return s.err
}

func getStack(skip int) string {
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
