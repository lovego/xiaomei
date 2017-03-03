package number

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
	switch value.(type) {
	case string:
		number = value.(string)
	case int:
		number = strconv.Itoa(value.(int))
	case int64:
		number = strconv.FormatInt(value.(int64), 10)
	case float64:
		number = fmt.Sprintf(`%0.2f`, value.(float64))
	default:
		panic(fmt.Sprintf(`unsupported value type: %v`, reflect.TypeOf(value).Name()))
	}
	return thousandSep(number)
}

func thousandSep(number string) string {
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
