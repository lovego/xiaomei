package union

import (
	"reflect"
)

func Union(a, b interface{}, keyName string) [][2]interface{} {
	aV, bV, swapped := makeABValue(a, b)
	bM := makeBMap(bV, keyName)
	result := unionAB(aV, bM, keyName, swapped)
	if len(bM) > 0 {
		result = addBvRemainder(result, bV, bM, keyName, swapped)
	}
	return result
}

func unionAB(
	aV reflect.Value, bM map[string]interface{}, keyName string, swapped bool,
) [][2]interface{} {
	var result [][2]interface{}
	length := aV.Len()
	for i := 0; i < length; i++ {
		av := aV.Index(i)
		key := getKey(av, keyName)
		bv, ok := bM[key]
		if ok {
			delete(bM, key)
		}
		if swapped {
			result = append(result, [2]interface{}{bv, av.Interface()})
		} else {
			result = append(result, [2]interface{}{av.Interface(), bv})
		}
	}
	return result
}

func addBvRemainder(result [][2]interface{},
	bV reflect.Value, bM map[string]interface{}, keyName string, swapped bool,
) [][2]interface{} {
	length := bV.Len()
	for i := 0; i < length; i++ {
		bv := bV.Index(i)
		key := getKey(bv, keyName)
		if _, ok := bM[key]; ok {
			if swapped {
				result = append(result, [2]interface{}{bv.Interface(), nil})
			} else {
				result = append(result, [2]interface{}{nil, bv.Interface()})
			}
		}
	}
	return result
}

func makeABValue(a, b interface{}) (reflect.Value, reflect.Value, bool) {
	if a == nil {
		a = []struct{}{}
	}
	if b == nil {
		b = []struct{}{}
	}
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)
	if bV.Len() > aV.Len() {
		return aV, bV, false
	} else {
		return bV, aV, true
	}
}

func makeBMap(bV reflect.Value, keyName string) map[string]interface{} {
	bM := make(map[string]interface{})
	length := bV.Len()
	for i := 0; i < length; i++ {
		v := bV.Index(i)
		bM[getKey(v, keyName)] = v.Interface()
	}
	return bM
}

func getKey(v reflect.Value, name string) string {
	if v.Kind() == reflect.Struct {
		return v.FieldByName(name).String()
	} else {
		return v.MapIndex(reflect.ValueOf(name)).Interface().(string)
	}
}
