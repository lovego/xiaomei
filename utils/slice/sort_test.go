package slice

import (
	"reflect"
	"testing"
)

type ssCase [2][]interface{}

func TestSort(t *testing.T) {
	var nil_result = []interface{}{}
	var cases = []ssCase{
		ssCase{nil_result, nil_result},
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
	for _, test_case := range cases {
		got := test_case[0]
		Sort(got, `V`)
		expect := test_case[1]

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("expect: %v, got: %v\n", expect, got)
		}
	}
}
