package merge

import (
	"reflect"
	"strings"
)

func Merge(a, b interface{}) interface{} {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if bv.Kind() != reflect.Map {
		panic(`b must be a map.`)
	}
	return mergeMap2MapOrStruct(av, bv).Interface()
}

func mergeMap2MapOrStruct(a, b reflect.Value) (result reflect.Value) {
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	if b.Kind() == reflect.Interface {
		b = b.Elem()
	}
	switch a.Kind() {
	case reflect.Map:
		result = shallowCopyMap(a)
		for _, key := range b.MapKeys() {
			if key.Kind() == reflect.Interface {
				key = key.Elem()
			}
			result.SetMapIndex(key, mergeMap2MapOrStruct(a.MapIndex(key), b.MapIndex(key)))
		}
	case reflect.Struct:
		result = shallowCopyStruct(a)
		for _, key := range b.MapKeys() {
			if key.Kind() == reflect.Interface {
				key = key.Elem()
			}
			keyStr := strings.Title(key.String())
			result.FieldByName(keyStr).Set(mergeMap2MapOrStruct(a.FieldByName(keyStr), b.MapIndex(key)))
		}
	case reflect.Slice:
		result = shallowCopySlice(a)
		for i, l := 0, b.Len(); i < l; i++ {
			v := b.Index(i)
			if v.Kind() == reflect.Interface {
				v = v.Elem()
			}
			result = reflect.Append(result, v)
		}
	default:
		return b
	}
	return result
}

// 浅拷贝，所以返回值并不能保证和入参完全独立
func shallowCopyMap(v reflect.Value) reflect.Value {
	dup := reflect.MakeMap(v.Type())
	for _, key := range v.MapKeys() {
		dup.SetMapIndex(key, v.MapIndex(key))
	}
	return dup
}

// 浅拷贝，所以返回值并不能保证和入参完全独立
func shallowCopyStruct(v reflect.Value) reflect.Value {
	dup := reflect.New(v.Type()).Elem()
	for i, l := 0, v.NumField(); i < l; i++ {
		dup.Field(i).Set(v.Field(i))
	}
	return dup
}

// 浅拷贝，所以返回值并不能保证和入参完全独立
func shallowCopySlice(v reflect.Value) reflect.Value {
	l := v.Len()
	dup := reflect.MakeSlice(v.Type(), l, l)
	for i := 0; i < l; i++ {
		dup.Index(i).Set(v.Index(i))
	}
	return dup
}
