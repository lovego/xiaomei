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

func TestMerge1(t *testing.T) {
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

func TestMerge2(t *testing.T) {
	type Service struct {
		Image, Ports     string
		Command, Volumes []string
	}
	config := struct {
		Services        map[string]Service
		VolumesToCreate []string `yaml:"volumesToCreate"`
		Environments    map[string]map[string]interface{}
	}{
		Services: map[string]Service{`app`: {
			Image:   `example/app`,
			Ports:   `3001~3003`,
			Volumes: []string{`example_logs:/home/ubuntu/example/log`},
		}},
		VolumesToCreate: []string{`example_logs`},
		Environments: map[string]map[string]interface{}{
			`dev`: map[string]interface{}{
				`services`: map[interface{}]interface{}{
					`app`: map[interface{}]interface{}{
						`ports`: `3001`,
					},
				},
			},
		},
	}
	newConfig := Merge(config, config.Environments[`dev`])
	fmt.Printf("%+v\n", newConfig)
}
