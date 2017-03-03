package number

import (
	"testing"
)

func TestThousand(t *testing.T) {
	cases := map[string]interface{}{
		`1`:            `1`,
		`12`:           12,
		`123`:          `123`,
		`1,234`:        `1234`,
		`12,345`:       `12345`,
		`123,456`:      `123456`,
		`1,234,567`:    `1234567`,
		`1,234,567.8`:  `1234567.8`,
		`1,234,567.89`: 1234567.89,
	}
	for expect, input := range cases {
		if got := ThousandSep(input); got != expect {
			t.Errorf("input: %s, expect: %s, got: %s\n", input, expect, got)
		}
	}
}
