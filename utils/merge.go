package utils

import (
	"reflect"
)

func Merge(ai, bi interface{}) interface{} {
	return merge(reflect.ValueOf(ai), reflect.ValueOf(ai)).Interface()
}

func merge(a, b reflect.Value) reflect.Value {
	switch b.Kind() {
	case reflect.Map:
		for _, key := range b.MapKeys() {
			if value := merge(a.MapIndex(key), b.MapIndex(key)); value.IsValid() {
				a.SetMapIndex(key, value)
			}
		}
		return b
	case reflect.Struct:
		return b
	default:
		return b
	}

}
