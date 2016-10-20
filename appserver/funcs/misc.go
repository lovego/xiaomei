package funcs

import (
	"reflect"
)

func IF(v, a, b interface{}) interface{} {
	if v == nil {
		return b
	}
	value := reflect.ValueOf(v)
	switch value.Type().Kind() {
	case reflect.Bool:
		if v.(bool) == false {
			return b
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		if value.Len() == 0 {
			return b
		}
	}
	return a
}
