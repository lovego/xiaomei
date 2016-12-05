package utils

import (
	"io"
	"os"
)

func Exist(p string) bool {
	fi, _ := os.Stat(p)
	return fi != nil
}

func IsDir(p string) bool {
	fi, _ := os.Stat(p)
	return fi != nil && fi.IsDir()
}

func IsEmptyDir(p string) bool {
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == io.EOF
}
