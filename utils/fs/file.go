package fs

import (
	"io"
	"os"
)

func OpenAppend(p string) (*os.File, error) {
	return os.OpenFile(p, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}

func Copy(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	info, err := os.Stat(src)
	if err != nil {
		return
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	return
}
