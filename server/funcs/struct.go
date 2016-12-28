package funcs

import (
	"reflect"
)

func StructOrMapField(obj interface{}, field string) interface{} {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		} else {
			value = reflect.Indirect(value)
		}
	}
	if value.Kind() == reflect.Map {
		v := value.MapIndex(reflect.ValueOf(field))
		if v.IsValid() {
			return v.Interface()
		} else {
			return nil
		}
	} else {
		if value.IsValid() {
			return value.FieldByName(field).Interface()
		} else {
			return nil
		}
	}
}
