package utils

import (
	"reflect"
)

func Merge(a, b interface{}) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if bv.Kind() != reflect.Map {
		panic(`b must be a map.`)
	}
	mergeMap(av, bv, ``)
}

func mergeMap(a, b reflect.Value, keys string) {
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	for _, key := range b.MapKeys() {
		bv := b.MapIndex(key)
		if bv.Kind() == reflect.Interface {
			bv = bv.Elem()
		}
		if bv.Kind() == reflect.Map {
			switch a.Kind() {
			case reflect.Map:
				mergeMap(a.MapIndex(key), bv, keys+`.`+key.String())
			case reflect.Struct:
				mergeMap(a.FieldByName(key.String()), bv, keys+`.`+key.String())
			default:
				panic(`a` + keys + ` should be a map or struct`)
			}
		} else {
			switch a.Kind() {
			case reflect.Map:
				a.SetMapIndex(key, bv)
			case reflect.Struct:
				a.FieldByName(key.String()).Set(bv)
			default:
				panic(`a` + keys + ` should be a map or struct`)
			}
		}
	}
}
