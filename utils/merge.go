package utils

import (
	"reflect"
)

func Merge(a, b interface{}) interface{} {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if bv.Kind() != reflect.Map {
		panic(`b must be a map.`)
	}
	return mergeMap2MapOrStruct(av, bv, ``).Interface()
}

func mergeMap2MapOrStruct(a, b reflect.Value, keys string) (result reflect.Value) {
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	if b.Kind() == reflect.Interface {
		b = b.Elem()
	}
	switch a.Kind() {
	case reflect.Map:
		result = reflect.MakeMap(a.Type())
		for _, key := range a.MapKeys() {
			result.SetMapIndex(key, mergeValue(a.MapIndex(key), b, key, keys))
		}
	case reflect.Struct:
		typ := a.Type()
		result = reflect.New(typ).Elem()
		for i, l := 0, a.NumField(); i < l; i++ {
			key := typ.Field(i).Name
			if len(key) > 1 {
				key = strings.ToLower(key[:1]) + key[1:]
			} else {
				key = strings.ToLower(key)
			}
			result.Field(i).Set(mergeValue(a.Field(i), b, reflect.ValueOf(key), keys))
		}
	default:
		return b
	}
	return result
}

func mergeValue(av, b, key reflect.Value, keys string) reflect.Value {
	if bv := b.MapIndex(key); bv.IsValid() {
		return mergeMap2MapOrStruct(av, bv, keys)
	} else {
		return av
	}
}
