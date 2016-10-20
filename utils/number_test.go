package utils

import (
	"gopkg.in/inf.v0"
	"strconv"
	"testing"
)

func TestThousandSep(t *testing.T) {
	cases := map[string]string{
		`1`:          `1`,
		`12`:         `12`,
		`123`:        `123`,
		`1234`:       `1,234`,
		`12345`:      `12,345`,
		`123456`:     `123,456`,
		`1234567`:    `1,234,567`,
		`1234567.8`:  `1,234,567.8`,
		`1234567.89`: `1,234,567.89`,
	}
	for input, expect := range cases {
		if got := ThousandSep(input); got != expect {
			t.Errorf("input: %s, expect: %s, got: %s\n", input, expect, got)
		}
		input = `-` + input
		expect = `-` + expect
		if got := ThousandSep(input); got != expect {
			t.Errorf("input: %s, expect: %s, got: %s\n", input, expect, got)
		}
	}
}

func TestThousands(t *testing.T) {
	cases := map[string]string{
		`4`:          `0.00`,
		`12`:         `0.01`,
		`123`:        `0.12`,
		`1234`:       `1.23`,
		`12345`:      `12.35`,
		`123456`:     `123.46`,
		`1234567`:    `1,234.57`,
		`1234567.8`:  `1,234.57`,
		`1234567.89`: `1,234.57`,
	}
	dec := new(inf.Dec)
	for input, expect := range cases {
		dec.SetString(input)
		if got := Thousands(dec); got != expect {
			t.Errorf("input: %s, expect: %s, got: %s\n", input, expect, got)
		}
		if i, _ := strconv.ParseFloat(input, 64); i >= 5 {
			expect = `-` + expect
		}
		input = `-` + input
		dec.SetString(input)
		if got := Thousands(dec); got != expect {
			t.Errorf("input: %s, expect: %s, got: %s\n", input, expect, got)
		}
	}
}
