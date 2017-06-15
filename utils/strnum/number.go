package strnum

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ThousandSep(value interface{}) string {
	if value == nil {
		return `0`
	}
	var number string
	switch v := value.(type) {
	case string:
		number = v
	case int:
		number = strconv.Itoa(v)
	case int64:
		number = strconv.FormatInt(v, 10)
	case float64:
		number = fmt.Sprintf(`%0.2f`, v)
	default:
		panic(fmt.Sprintf(`unsupported value type: %v`, reflect.TypeOf(value).Name()))
	}
	return thousandSep(number)
}

func thousandSep(number string) string {
	if number == `` {
		return `0`
	}
	start, end := 0, len(number)
	if number[0] == '-' {
		start = 1
	}
	if i := strings.IndexByte(number, '.'); i >= 0 {
		end = i
	}
	return number[:start] +
		strings.Join(ThousandSplit(number[start:end]), `,`) +
		number[end:]
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
