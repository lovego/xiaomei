package union

import (
	"reflect"
	"testing"
)

type susCase [2][][2]interface{}

func TestUnionSort(t *testing.T) {
	var nilResult [][2]interface{}
	var cases = []susCase{
		susCase{nilResult, nilResult},
		susCase{
			[][2]interface{}{
				[2]interface{}{s{`3`, 3}},
				[2]interface{}{s{`1`, 1}},
				[2]interface{}{s{`2`, 2}},
			},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}},
				[2]interface{}{s{`2`, 2}},
				[2]interface{}{s{`3`, 3}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{nil, s{`3`, 3}},
				[2]interface{}{nil, s{`1`, 1}},
				[2]interface{}{nil, s{`2`, 2}},
			},
			[][2]interface{}{
				[2]interface{}{nil, s{`1`, 1}},
				[2]interface{}{nil, s{`2`, 2}},
				[2]interface{}{nil, s{`3`, 3}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
			},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{s{`4`, 4}},
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
			},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{s{`4`, 4}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{nil, s{`4`, 4}},
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
			},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{nil, s{`4`, 4}},
			},
		},

		susCase{
			[][2]interface{}{
				[2]interface{}{m{`Name`: `3`, `V`: 3}},
				[2]interface{}{m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}},
			},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{nil, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{nil, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{nil, m{`Name`: `2`, `V`: 2}},
			},
			[][2]interface{}{
				[2]interface{}{nil, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{nil, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{nil, m{`Name`: `3`, `V`: 3}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
			},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{m{`Name`: `4`, `V`: 4}},
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
			},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{m{`Name`: `4`, `V`: 4}},
			},
		},
		susCase{
			[][2]interface{}{
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{nil, m{`Name`: `4`, `V`: 4}},
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
			},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{nil, m{`Name`: `4`, `V`: 4}},
			},
		},
	}
	for _, testCase := range cases {
		got := testCase[0]
		UnionSort(got, `V`)
		expect := testCase[1]

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("expect: %v, got: %v\n", expect, got)
		}
	}
}
