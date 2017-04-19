package utils

import (
	"fmt"
	"reflect"
	"syscall"
	"testing"
)

func TestMaximizeNOFILE(t *testing.T) {
	printNOFILE()
	MaximizeNOFILE()
	printNOFILE()
}

func printNOFILE() {
	rlimit := syscall.Rlimit{}
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		println(`get RLIMIT_NOFILE error: `, err.Error())
	} else {
		fmt.Println(rlimit)
	}
}

func TestStack(t *testing.T) {
	fmt.Printf("%s\n", Stack(1))
}

func TestReflect(t *testing.T) {
	v := reflect.ValueOf(``)
	fmt.Println(v.IsValid(), v)
}
