package fs

import (
	"os"
)

func OpenAppend(p string) *os.File {
	if f, err := os.OpenFile(p, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		panic(err)
	} else {
		return f
	}
}
