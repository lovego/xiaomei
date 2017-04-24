package utils

import (
	"reflect"
	"strings"
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
		if key.Kind() == reflect.Interface {
			key = key.Elem()
		}
		if bv.Kind() == reflect.Map {
			switch a.Kind() {
			case reflect.Map:
				mergeMap(a.MapIndex(key), bv, keys+`.`+key.String())
			case reflect.Struct:
				keyStr := strings.Title(key.String())
				mergeMap(a.FieldByName(keyStr), bv, keys+`.`+keyStr)
			default:
				panic(`a` + keys + ` should be a map or struct`)
			}
		} else {
			switch a.Kind() {
			case reflect.Map:
				if key.Kind() == reflect.Interface {
					key = key.Elem()
				}
				a.SetMapIndex(key, bv)
			case reflect.Struct:
				keyStr := strings.Title(key.String())
				if av := a.FieldByName(keyStr); av.IsValid() {
					println(keyStr)
					av.Set(bv)
				} else {
					panic(`no such field: ` + keyStr)
				}
			default:
				panic(`a` + keys + ` should be a map or struct`)
			}
		}
	}
}
