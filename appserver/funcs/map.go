package funcs

import (
	"errors"
	"reflect"
)

func MakeDict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func MapKeys(m interface{}) []interface{} {
	keys := reflect.ValueOf(m).MapKeys()
	result := make([]interface{}, len(keys))
	for i, v := range keys {
		result[i] = v.Interface()
	}
	return result
}

func MapValues(m interface{}) []interface{} {
	value := reflect.ValueOf(m)
	keys := value.MapKeys()
	values := make([]interface{}, len(keys))
	for i, key := range keys {
		values[i] = value.MapIndex(key).Interface()
	}
	return values
}

func MapKeysUnion(a, b interface{}) (result []string) {
	if a != nil {
		for _, value := range reflect.ValueOf(a).MapKeys() {
			result = append(result, value.Interface().(string))
		}
	}
	if b != nil {
	loop:
		for _, value := range reflect.ValueOf(b).MapKeys() {
			k := value.Interface().(string)
			for _, v := range result {
				if v == k {
					continue loop
				}
			}
			result = append(result, k)
		}
	}
	return
}
