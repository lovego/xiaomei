package merge

import (
	"fmt"
	"testing"
)

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
			Volumes: []string{`-v`, `example_logs:/home/ubuntu/example/log`},
		}},
		VolumesToCreate: []string{`example_logs`},
		Environments: map[string]map[string]interface{}{
			`dev`: map[string]interface{}{
				`services`: map[interface{}]interface{}{
					`app`: map[interface{}]interface{}{
						`ports`:   `3001`,
						`volumes`: []interface{}{`--appendix`},
					},
				},
			},
		},
	}
	newConfig := Merge(config, config.Environments[`dev`])
	fmt.Printf("%+v\n", newConfig)
}
