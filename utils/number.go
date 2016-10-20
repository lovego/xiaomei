package utils

import (
	"gopkg.in/inf.v0"
	"strings"
)

func ThousandSep(dec *inf.Dec) string {
	if dec == nil {
		return `0`
	}
	number := dec.String()
	start, end := 0, len(number)
	if number[0] == '-' {
		start = 1
	}
	if i := strings.IndexByte(number, '.'); i >= 0 {
		end = i
	}
	return number[:start] +
		strings.Join(ThousandSplit(number[start:end]), `,`) +
		number[end:len(number)]
}

func ThousandSplit(number string) (result []string) {
	start := 0
	if rem := len(number) % 3; rem > 0 {
		result = append(result, number[:rem])
		start = rem
	}
	for ; start < len(number); start += 3 {
		result = append(result, number[start:start+3])
	}
	return
}
