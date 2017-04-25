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

func TestMerge(t *testing.T) {
	type s struct {
		E int
		F string
	}
	a := struct {
		A, B string
		C, D int
		S    s
	}{A: `A`, C: 3, S: s{E: 5}}

	b := map[interface{}]interface{}{
		`B`: `B`, `D`: 4,
		`S`: map[interface{}]interface{}{`F`: `F`},
	}

	m := Merge(a, b)
	fmt.Printf("%+v, %+v, %+v\n", a, b, m)
}
