package slice

import (
	"reflect"
	"testing"
)

type ssCase [2][]interface{}

type s struct {
	Name string
	V    int
}

type m map[string]interface{}

func TestSort(t *testing.T) {
	var nilResult = []interface{}{}
	var cases = []ssCase{
		ssCase{nilResult, nilResult},
		ssCase{
			[]interface{}{
				s{`3`, 3},
				s{`1`, 1},
				s{`2`, 2},
			},
			[]interface{}{
				s{`1`, 1},
				s{`2`, 2},
				s{`3`, 3},
			},
		},

		ssCase{
			[]interface{}{
				m{`Name`: `3`, `V`: 3},
				m{`Name`: `1`, `V`: 1},
				m{`Name`: `2`, `V`: 2},
			},
			[]interface{}{
				m{`Name`: `1`, `V`: 1},
				m{`Name`: `2`, `V`: 2},
				m{`Name`: `3`, `V`: 3},
			},
		},
	}
	for _, testCase := range cases {
		got := testCase[0]
		Sort(got, `V`)
		expect := testCase[1]

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("expect: %v, got: %v\n", expect, got)
		}
	}
}
