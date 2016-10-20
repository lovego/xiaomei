package utils

import (
	"reflect"
	"strings"

	"gopkg.in/inf.v0"
)

func Escape(i interface{}, data []map[string]interface{}) interface{} {
	k := reflect.TypeOf(i).Kind()
	switch k {
	case reflect.Slice:
		value := reflect.ValueOf(i)
		for _, row := range data {
			val := reflect.New(reflect.TypeOf(i).Elem())
			setValue(val, row)
			value = reflect.Append(value, val.Elem())
		}
		return value.Interface()
	case reflect.Ptr:
		value := reflect.ValueOf(i)
		for _, row := range data {
			setValue(value, row)
		}
		return value.Elem().Interface()
	}
	return nil
}

func setValue(v reflect.Value, row map[string]interface{}) {
	s := v.Elem()
	for i := 0; i < s.NumField(); i++ {
		key := reflect.TypeOf(s.Interface()).Field(i).Name
		for k, val := range row {
			if strings.Contains(k, `_`) {
				k = strings.Replace(k, `_`, ``, -1)
			}
			if strings.ToLower(key) == k && s.Field(i).CanSet() {
				assert(s.Field(i), val)
			}
		}
	}
}

func assert(field reflect.Value, v interface{}) {
	switch field.Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(v.(string)))
	case reflect.Int:
		field.Set(reflect.ValueOf(v.(int)))
	case reflect.Float32:
		field.Set(reflect.ValueOf(v.(float32)))
	case reflect.Float64:
		field.Set(reflect.ValueOf(v.(float64)))
	case reflect.Ptr:
		field.Set(reflect.ValueOf(v.(*inf.Dec)))
	}
}
