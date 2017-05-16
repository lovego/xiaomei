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
	return merge(av, bv).Interface()
}

func merge(a, b reflect.Value) reflect.Value {
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	if b.Kind() == reflect.Interface {
		b = b.Elem()
	}
	switch a.Kind() {
	case reflect.Map:
		return mergeMap2Map(a, b)
	case reflect.Struct:
		return mergeMap2Struct(a, b)
	case reflect.Slice:
		return mergeSlice2Slice(a, b)
	default:
		return b
	}
}

func mergeMap2Map(a, b reflect.Value) reflect.Value {
	result := shallowCopyMap(a)
	for _, key := range b.MapKeys() {
		if key.Kind() == reflect.Interface {
			key = key.Elem()
		}
		result.SetMapIndex(key, merge(a.MapIndex(key), b.MapIndex(key)))
	}
	return result
}

func mergeMap2Struct(a, b reflect.Value) reflect.Value {
	result := shallowCopyStruct(a)
	for _, key := range b.MapKeys() {
		if key.Kind() == reflect.Interface {
			key = key.Elem()
		}
		keyStr := strings.Title(key.String())
		result.FieldByName(keyStr).Set(merge(a.FieldByName(keyStr), b.MapIndex(key)))
	}
	return result
}

func mergeSlice2Slice(a, b reflect.Value) reflect.Value {
	result := shallowCopySlice(a)
	for i, l := 0, b.Len(); i < l; i++ {
		v := b.Index(i)
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		result = reflect.Append(result, v)
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

func shallowCopyStruct(v reflect.Value) reflect.Value {
	dup := reflect.New(v.Type()).Elem()
	for i, l := 0, v.NumField(); i < l; i++ {
		dup.Field(i).Set(v.Field(i))
	}
	return dup
}

func shallowCopySlice(v reflect.Value) reflect.Value {
	l := v.Len()
	dup := reflect.MakeSlice(v.Type(), l, l)
	for i := 0; i < l; i++ {
		dup.Index(i).Set(v.Index(i))
	}
	return dup
}
