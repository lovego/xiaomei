package utils

import (
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
)

func MaximizeNOFILE() {
	bs, err := ioutil.ReadFile(`/proc/sys/fs/file-max`)
	if err != nil {
		println(`get file-max error: `, err.Error())
	}

	fileMax, err := strconv.Atoi(strings.TrimSpace(string(bs)))
	if err != nil {
		println(`atoi file-max error: `, err.Error())
	}

	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{
		Cur: uint64(fileMax), Max: uint64(fileMax),
	}); err != nil {
		println(`set RLIMIT_NOFILE error: `, err.Error())
	}
}
