package slice

import (
	"reflect"
	"testing"
)

type suCase [3]interface{}

type s struct {
	Name string
	V    int
}

type m map[string]interface{}

func TestUnion(t *testing.T) {
	var nil_result [][2]interface{}
	var cases = []suCase{
		suCase{nil, nil, nil_result},
		suCase{[]int{}, nil, nil_result},
		suCase{nil, []bool{}, nil_result},
		suCase{[]int{}, []bool{}, nil_result},
		suCase{
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}},
			[]s{},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}},
				[2]interface{}{s{`2`, 2}},
				[2]interface{}{s{`3`, 3}},
			},
		},
		suCase{
			[]s{},
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}},
			[][2]interface{}{
				[2]interface{}{nil, s{`1`, 1}},
				[2]interface{}{nil, s{`2`, 2}},
				[2]interface{}{nil, s{`3`, 3}},
			},
		},
		suCase{
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}},
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
			},
		},
		suCase{
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}, s{`4`, 4}},
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{s{`4`, 4}},
			},
		},
		suCase{
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}},
			[]s{s{`1`, 1}, s{`2`, 2}, s{`3`, 3}, s{`4`, 4}},
			[][2]interface{}{
				[2]interface{}{s{`1`, 1}, s{`1`, 1}},
				[2]interface{}{s{`2`, 2}, s{`2`, 2}},
				[2]interface{}{s{`3`, 3}, s{`3`, 3}},
				[2]interface{}{nil, s{`4`, 4}},
			},
		},
		suCase{
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}},
			[]m{},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}},
			},
		},
		suCase{
			[]m{},
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}},
			[][2]interface{}{
				[2]interface{}{nil, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{nil, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{nil, m{`Name`: `3`, `V`: 3}},
			},
		},
		suCase{
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}},
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
			},
		},
		suCase{
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}, m{`Name`: `4`, `V`: 4}},
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{m{`Name`: `4`, `V`: 4}},
			},
		},
		suCase{
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}},
			[]m{m{`Name`: `1`, `V`: 1}, m{`Name`: `2`, `V`: 2}, m{`Name`: `3`, `V`: 3}, m{`Name`: `4`, `V`: 4}},
			[][2]interface{}{
				[2]interface{}{m{`Name`: `1`, `V`: 1}, m{`Name`: `1`, `V`: 1}},
				[2]interface{}{m{`Name`: `2`, `V`: 2}, m{`Name`: `2`, `V`: 2}},
				[2]interface{}{m{`Name`: `3`, `V`: 3}, m{`Name`: `3`, `V`: 3}},
				[2]interface{}{nil, m{`Name`: `4`, `V`: 4}},
			},
		},
	}
	for _, test_case := range cases {
		got := Union(test_case[0], test_case[1], `Name`)
		expect := test_case[2].([][2]interface{})

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("input: %v, %v, expect: %v, got: %v\n", test_case[0], test_case[1], expect, got)
		}
	}
}
