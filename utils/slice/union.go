package slice

import (
	"reflect"
)

func Union(a, b interface{}, key_name string) (result [][2]interface{}) {
	if a == nil {
		a = []struct{}{}
	}
	if b == nil {
		b = []struct{}{}
	}
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	swapped := false
	if bv.Len() > av.Len() {
		av, bv = bv, av
		swapped = true
	}

	bm := make(map[string]interface{})
	length := bv.Len()
	for i := 0; i < length; i++ {
		v := bv.Index(i)
		bm[getKey(v, key_name)] = v.Interface()
	}

	length = av.Len()
	for i := 0; i < length; i++ {
		a_v := av.Index(i)
		key := getKey(a_v, key_name)
		b_v, ok := bm[key]
		if ok {
			delete(bm, key)
		}
		if swapped {
			result = append(result, [2]interface{}{b_v, a_v.Interface()})
		} else {
			result = append(result, [2]interface{}{a_v.Interface(), b_v})
		}
	}

	for _, b_v := range bm {
		if swapped {
			result = append(result, [2]interface{}{b_v, nil})
		} else {
			result = append(result, [2]interface{}{nil, b_v})
		}
	}

	return
}

func getKey(v reflect.Value, name string) string {
	if v.Kind() == reflect.Struct {
		return v.FieldByName(name).String()
	} else {
		return v.MapIndex(reflect.ValueOf(name)).Interface().(string)
	}
}
